package wikipedia

import (
	"errors"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

var dbEndpoint = "http://localhost:17474"

func TestIsValidCrawlLink(t *testing.T) {
	t.Run("does not crawl on links with ':'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("/wiki/Category:Spinash"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki/Test:"), false)
	})
	t.Run("does not crawl on links not starting with '/wiki/'", func(t *testing.T) {
		assert.Equal(t, IsValidCrawlLink("https://wikipedia.org"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki"), false)
		assert.Equal(t, IsValidCrawlLink("wikipedia/wiki/"), false)
		assert.Equal(t, IsValidCrawlLink("/wiki/binary"), true)
	})
}

func TestAddToDb(t *testing.T) {
	// keep errors in array
	errors := []string{}
	logErr = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}
	t.Run("fails when no server found", func(t *testing.T) {
		os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
		// first test bad response
		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet", []string{"/wiki/Animal"})
		assert.EqualError(t, err, "Post http://localhost:17474/edges?node=25079: dial tcp 127.0.0.1:17474: connect: connection refused")
		assert.Equal(t, []string{}, newNodes)
	})
	t.Run("returns error when current node doesnt exist (404)", func(t *testing.T) {
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", "https://en.wikipedia.org/wiki/Pet_door",
			httpmock.NewStringResponder(200, PetDoorResponse))
		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(404, map[string]interface{}{
					"code":  404,
					"error": "Node was not found",
				})
			},
		)

		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal"})
		assert.EqualError(t, err, "Node was not found")
		assert.Equal(t, newNodes, []string{})
	})
	t.Run("succesfully adds neighbor nodes", func(t *testing.T) {
		errors = []string{}
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", "https://en.wikipedia.org/wiki/Pet_door",
			httpmock.NewStringResponder(200, PetDoorResponse))
		httpmock.RegisterResponder("GET", "https://en.wikipedia.org/wiki/Animal",
			httpmock.NewStringResponder(200, animalResponse))

		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{11039790}})
			},
		)

		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal"})
		assert.Nil(t, err)
		assert.Equal(t, errors, []string{})
		assert.Equal(t, newNodes, []string{"https://en.wikipedia.org/wiki/Animal"})
	})
	t.Run("only returns new neighbors", func(t *testing.T) {
		// mock out http endpoint
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", "https://en.wikipedia.org/wiki/Pet_door",
			httpmock.NewStringResponder(200, PetDoorResponse))
		httpmock.RegisterResponder("GET", "https://en.wikipedia.org/wiki/Animal",
			httpmock.NewStringResponder(200, animalResponse))
		httpmock.RegisterResponder("GET", "https://en.wikipedia.org/wiki/Petula_Clark",
			httpmock.NewStringResponder(200, PetulaClarkResponse))

		httpmock.RegisterResponder("POST", dbEndpoint+"/edges?node=3276454",
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(200, map[string]interface{}{"neighborsAdded": []int{11039790}})
			},
		)

		newNodes, err := AddEdgesIfDoNotExist("/wiki/Pet_door", []string{"/wiki/Animal", "/wiki/Petula_clark"})
		assert.Nil(t, err)
		assert.Equal(t, newNodes, []string{"https://en.wikipedia.org/wiki/Animal"})
	})
}

func TestConnectToDB(t *testing.T) {
	dbEndpoint := "http://localhost:17474"
	os.Setenv("GRAPH_DB_ENDPOINT", dbEndpoint)
	t.Run("fails when db not found", func(t *testing.T) {
		err := ConnectToDB()
		assert.EqualError(t, err, "Get http://localhost:17474/metrics: dial tcp 127.0.0.1:17474: connect: connection refused")
	})
	t.Run("succeed when server exists", func(t *testing.T) {
		// mock out http endpoint
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", dbEndpoint+"/metrics",
			httpmock.NewStringResponder(200, `TEST`))

		err := ConnectToDB()
		assert.Nil(t, err)
	})
}

func TestGetArticleId(t *testing.T) {
	os.Setenv("WIKI_API_ENDPOINT", "https://en.wikipedia.org/w/api.php")
	t.Run("returns error on bad url", func(t *testing.T) {
		id, err := getArticleId("/wiki/DFSDfet_doorSDFUSFU#UFFISd")
		assert.NotNil(t, err)
		assert.Equal(t, -1, id)
	})
	t.Run("returns correct values", func(t *testing.T) {
		type test struct {
			Page          string
			ExpectedValue int
			ExpectedError error
		}

		table := []test{
			test{
				"/wiki/Pet",
				25079,
				nil,
			},
			test{
				"/wiki/Animal",
				11039790,
				nil,
			},
			test{
				"/wiki/Tests_(album)",
				24088322,
				nil,
			},
			test{
				"/wiki/The_Microphones",
				847580,
				nil,
			},
			test{
				"/wiki/SDF32fj302jf",
				-1,
				errors.New("status code error: 404 404 Not Found"),
			},
		}

		for _, v := range table {
			t.Run(v.Page, func(t *testing.T) {
				id, err := getArticleId(v.Page)
				assert.Equal(t, id, v.ExpectedValue)
				assert.Equal(t, err, v.ExpectedError)
			})
		}
	})

}

var PetDoorResponse = `
<!DOCTYPE html>
<html class="client-nojs" lang="en" dir="ltr">
<head>
<meta charset="UTF-8"/>
<title>Pet door - Wikipedia</title>
<script>document.documentElement.className=document.documentElement.className.replace(/(^|\s)client-nojs(\s|$)/,"$1client-js$2");RLCONF={"wgCanonicalNamespace":"","wgCanonicalSpecialPageName":!1,"wgNamespaceNumber":0,"wgPageName":"Pet_door","wgTitle":"Pet door","wgCurRevisionId":874746883,"wgRevisionId":874746883,"wgArticleId":3276454,"wgIsArticle":!0,"wgIsRedirect":!1,"wgAction":"view","wgUserName":null,"wgUserGroups":["*"],"wgCategories":["Wikipedia articles needing clarification from September 2012","Articles containing Spanish-language text","Articles containing French-language text","Articles containing Middle English-language text","Articles needing additional references from July 2007","All articles needing additional references","Commons category link is on Wikidata","Door furniture","Pet equipment"],"wgBreakFrames":!1,"wgPageContentLanguage":"en","wgPageContentModel":"wikitext","wgSeparatorTransformTable":["",""],"wgDigitTransformTable":["",""],"wgDefaultDateFormat":"dmy",
"wgMonthNames":["","January","February","March","April","May","June","July","August","September","October","November","December"],"wgMonthNamesShort":["","Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"],"wgRelevantPageName":"Pet_door","wgRelevantArticleId":3276454,"wgRequestId":"XS9QawpAICAAAHE1598AAABD","wgCSPNonce":!1,"wgIsProbablyEditable":!0,"wgRelevantPageIsProbablyEditable":!0,"wgRestrictionEdit":[],"wgRestrictionMove":[],"wgMediaViewerOnClick":!0,"wgMediaViewerEnabledByDefault":!0,"wgPopupsReferencePreviews":!1,"wgPopupsConflictsWithNavPopupGadget":!1,"wgVisualEditor":{"pageLanguageCode":"en","pageLanguageDir":"ltr","pageVariantFallbacks":"en"},"wgMFDisplayWikibaseDescriptions":{"search":!0,"nearby":!0,"watchlist":!0,"tagline":!1},"wgWMESchemaEditAttemptStepOversample":!1,"wgPoweredByHHVM":!0,"wgULSCurrentAutonym":"English","wgNoticeProject":"wikipedia","wgWikibaseItemId":"Q943110","wgCentralAuthMobileDomain":!1,
"wgEditSubmitButtonLabelPublish":!0};RLSTATE={"ext.gadget.charinsert-styles":"ready","ext.globalCssJs.user.styles":"ready","ext.globalCssJs.site.styles":"ready","site.styles":"ready","noscript":"ready","user.styles":"ready","ext.globalCssJs.user":"ready","ext.globalCssJs.site":"ready","user":"ready","user.options":"ready","user.tokens":"loading","ext.cite.styles":"ready","mediawiki.legacy.shared":"ready","mediawiki.legacy.commonPrint":"ready","mediawiki.toc.styles":"ready","wikibase.client.init":"ready","ext.visualEditor.desktopArticleTarget.noscript":"ready","ext.uls.interlanguage":"ready","ext.wikimediaBadges":"ready","ext.3d.styles":"ready","mediawiki.skinning.interface":"ready","skins.vector.styles":"ready"};RLPAGEMODULES=["ext.cite.ux-enhancements","site","mediawiki.page.startup","mediawiki.page.ready","mediawiki.toc","mediawiki.searchSuggest","ext.gadget.teahouse","ext.gadget.ReferenceTooltips","ext.gadget.watchlist-notice","ext.gadget.DRN-wizard","ext.gadget.charinsert",
"ext.gadget.refToolbar","ext.gadget.extra-toolbar-buttons","ext.gadget.switcher","ext.centralauth.centralautologin","mmv.head","mmv.bootstrap.autostart","ext.popups","ext.visualEditor.desktopArticleTarget.init","ext.visualEditor.targetLoader","ext.eventLogging","ext.wikimediaEvents","ext.navigationTiming","ext.uls.compactlinks","ext.uls.interface","ext.quicksurveys.init","ext.centralNotice.geoIP","ext.centralNotice.startUp","skins.vector.js"];</script>
<script>(RLQ=window.RLQ||[]).push(function(){mw.loader.implement("user.tokens@0tffind",function($,jQuery,require,module){/*@nomin*/mw.user.tokens.set({"editToken":"+\\","patrolToken":"+\\","watchToken":"+\\","csrfToken":"+\\"});
});});</script>
</html>
`

var PetulaClarkResponse = `
<!DOCTYPE html>
<html class="client-nojs" lang="en" dir="ltr">
<head>
<meta charset="UTF-8"/>
<title>Petula Clark - Wikipedia</title>
<script>document.documentElement.className=document.documentElement.className.replace(/(^|\s)client-nojs(\s|$)/,"$1client-js$2");RLCONF={"wgCanonicalNamespace":"","wgCanonicalSpecialPageName":!1,"wgNamespaceNumber":0,"wgPageName":"Petula_Clark","wgTitle":"Petula Clark","wgCurRevisionId":904356293,"wgRevisionId":904356293,"wgArticleId":197772,"wgIsArticle":!0,"wgIsRedirect":!1,"wgAction":"view","wgUserName":null,"wgUserGroups":["*"],"wgCategories":["Use dmy dates from December 2017","Use British English from March 2012","BLP articles lacking sources from March 2012","All BLP articles lacking sources","Articles with hCards","All articles with unsourced statements","Articles with unsourced statements from February 2011","Articles with unsourced statements from October 2012","Articles with unsourced statements from July 2011","Commons category link is on Wikidata","Wikipedia articles with BIBSYS identifiers","Wikipedia articles with BNE identifiers",
"Wikipedia articles with BNF identifiers","Wikipedia articles with GND identifiers","Wikipedia articles with ISNI identifiers","Wikipedia articles with LCCN identifiers","Wikipedia articles with MusicBrainz identifiers","Wikipedia articles with NKC identifiers","Wikipedia articles with SNAC-ID identifiers","Wikipedia articles with SUDOC identifiers","Wikipedia articles with VIAF identifiers","Wikipedia articles with WorldCat-VIAF identifiers","1932 births","20th-century English actresses","20th-century English singers","21st-century English actresses","21st-century English singers","British Invasion artists","Commanders of the Order of the British Empire","English child actresses","English child singers","English film actresses","English film score composers","English musical theatre actresses","English musical theatre composers","English people of Welsh descent","English television actresses","English expatriates in Switzerland","French-language singers","German-language singers",
"Grammy Award winners","Italian-language singers","Living people","MGM Records artists","Musicians from Surrey","People from Epsom","Pye Records artists","English expatriates in France","Decca Records artists","EMI Records artists","Imperial Records artists","Warner Bros. Records artists","Columbia Records artists","Schlager musicians","English female pop singers","20th-century women singers","21st-century women singers"],"wgBreakFrames":!1,"wgPageContentLanguage":"en","wgPageContentModel":"wikitext","wgSeparatorTransformTable":["",""],"wgDigitTransformTable":["",""],"wgDefaultDateFormat":"dmy","wgMonthNames":["","January","February","March","April","May","June","July","August","September","October","November","December"],"wgMonthNamesShort":["","Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"],"wgRelevantPageName":"Petula_Clark","wgRelevantArticleId":197772,"wgRequestId":"XTCc6ApAMEwAAJ36s5AAAAAD","wgCSPNonce":!1,"wgIsProbablyEditable":!0,
"wgRelevantPageIsProbablyEditable":!0,"wgRestrictionEdit":[],"wgRestrictionMove":[],"wgMediaViewerOnClick":!0,"wgMediaViewerEnabledByDefault":!0,"wgPopupsReferencePreviews":!1,"wgPopupsConflictsWithNavPopupGadget":!1,"wgVisualEditor":{"pageLanguageCode":"en","pageLanguageDir":"ltr","pageVariantFallbacks":"en"},"wgMFDisplayWikibaseDescriptions":{"search":!0,"nearby":!0,"watchlist":!0,"tagline":!1},"wgWMESchemaEditAttemptStepOversample":!1,"wgPoweredByHHVM":!0,"wgULSCurrentAutonym":"English","wgNoticeProject":"wikipedia","wgWikibaseItemId":"Q236212","wgCentralAuthMobileDomain":!1,"wgEditSubmitButtonLabelPublish":!0};RLSTATE={"ext.gadget.charinsert-styles":"ready","ext.globalCssJs.user.styles":"ready","ext.globalCssJs.site.styles":"ready","site.styles":"ready","noscript":"ready","user.styles":"ready","ext.globalCssJs.user":"ready","ext.globalCssJs.site":"ready","user":"ready","user.options":"ready","user.tokens":"loading","ext.cite.styles":"ready",
"mediawiki.legacy.shared":"ready","mediawiki.legacy.commonPrint":"ready","mediawiki.toc.styles":"ready","wikibase.client.init":"ready","ext.visualEditor.desktopArticleTarget.noscript":"ready","ext.uls.interlanguage":"ready","ext.wikimediaBadges":"ready","ext.3d.styles":"ready","mediawiki.skinning.interface":"ready","skins.vector.styles":"ready"};RLPAGEMODULES=["ext.cite.ux-enhancements","site","mediawiki.page.startup","mediawiki.page.ready","mediawiki.toc","mediawiki.searchSuggest","ext.gadget.teahouse","ext.gadget.ReferenceTooltips","ext.gadget.watchlist-notice","ext.gadget.DRN-wizard","ext.gadget.charinsert","ext.gadget.refToolbar","ext.gadget.extra-toolbar-buttons","ext.gadget.switcher","ext.centralauth.centralautologin","mmv.head","mmv.bootstrap.autostart","ext.popups","ext.visualEditor.desktopArticleTarget.init","ext.visualEditor.targetLoader","ext.eventLogging","ext.wikimediaEvents","ext.navigationTiming","ext.uls.compactlinks","ext.uls.interface","ext.quicksurveys.init",
"ext.centralNotice.geoIP","ext.centralNotice.startUp","skins.vector.js"];</script>
<script>(RLQ=window.RLQ||[]).push(function(){mw.loader.implement("user.tokens@0tffind",function($,jQuery,require,module){/*@nomin*/mw.user.tokens.set({"editToken":"+\\","patrolToken":"+\\","watchToken":"+\\","csrfToken":"+\\"});
});});</script>
</html>
`

var animalResponse = `
<!DOCTYPE html>
<html class="client-nojs" lang="en" dir="ltr">
<head>
<meta charset="UTF-8"/>
<title>Animal - Wikipedia</title>
<script>document.documentElement.className=document.documentElement.className.replace(/(^|\s)client-nojs(\s|$)/,"$1client-js$2");RLCONF={"wgCanonicalNamespace":"","wgCanonicalSpecialPageName":!1,"wgNamespaceNumber":0,"wgPageName":"Animal","wgTitle":"Animal","wgCurRevisionId":906447531,"wgRevisionId":906447531,"wgArticleId":11039790,"wgIsArticle":!0,"wgIsRedirect":!1,"wgAction":"view","wgUserName":null,"wgUserGroups":["*"],"wgCategories":["CS1 maint: Untitled periodical","CS1: long volume value","CS1 Latin-language sources (la)","CS1 German-language sources (de)","Wikipedia indefinitely semi-protected pages","Wikipedia indefinitely move-protected pages","Articles with short description","Good articles","Use dmy dates from October 2012","Use British English from April 2017","Articles with 'species' microformats","Articles containing Latin-language text","Articles containing potentially dated statements from 2013","All articles containing potentially dated statements",
"Wikipedia articles with BNF identifiers","Wikipedia articles with GND identifiers","Wikipedia articles with LCCN identifiers","Wikipedia articles with NARA identifiers","Wikipedia articles with NDL identifiers","Animals","Kingdoms (biology)","Cryogenian first appearances","Taxa named by Carl Linnaeus"],"wgBreakFrames":!1,"wgPageContentLanguage":"en","wgPageContentModel":"wikitext","wgSeparatorTransformTable":["",""],"wgDigitTransformTable":["",""],"wgDefaultDateFormat":"dmy","wgMonthNames":["","January","February","March","April","May","June","July","August","September","October","November","December"],"wgMonthNamesShort":["","Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"],"wgRelevantPageName":"Animal","wgRelevantArticleId":11039790,"wgRequestId":"XTAVagpAAEQAABIO@XgAAAAL","wgCSPNonce":!1,"wgIsProbablyEditable":!1,"wgRelevantPageIsProbablyEditable":!1,"wgRestrictionEdit":["autoconfirmed"],"wgRestrictionMove":["sysop"],"wgMediaViewerOnClick":!0,
"wgMediaViewerEnabledByDefault":!0,"wgPopupsReferencePreviews":!1,"wgPopupsConflictsWithNavPopupGadget":!1,"wgVisualEditor":{"pageLanguageCode":"en","pageLanguageDir":"ltr","pageVariantFallbacks":"en"},"wgMFDisplayWikibaseDescriptions":{"search":!0,"nearby":!0,"watchlist":!0,"tagline":!1},"wgWMESchemaEditAttemptStepOversample":!1,"wgPoweredByHHVM":!0,"wgULSCurrentAutonym":"English","wgNoticeProject":"wikipedia","wgWikibaseItemId":"Q729","wgCentralAuthMobileDomain":!1,"wgEditSubmitButtonLabelPublish":!0};RLSTATE={"ext.gadget.charinsert-styles":"ready","ext.globalCssJs.user.styles":"ready","ext.globalCssJs.site.styles":"ready","site.styles":"ready","noscript":"ready","user.styles":"ready","ext.globalCssJs.user":"ready","ext.globalCssJs.site":"ready","user":"ready","user.options":"ready","user.tokens":"loading","ext.cite.styles":"ready","mediawiki.legacy.shared":"ready","mediawiki.legacy.commonPrint":"ready","jquery.makeCollapsible.styles":"ready",
"mediawiki.toc.styles":"ready","wikibase.client.init":"ready","ext.visualEditor.desktopArticleTarget.noscript":"ready","ext.uls.interlanguage":"ready","ext.wikimediaBadges":"ready","ext.3d.styles":"ready","mediawiki.skinning.interface":"ready","skins.vector.styles":"ready"};RLPAGEMODULES=["ext.cite.ux-enhancements","ext.scribunto.logs","site","mediawiki.page.startup","mediawiki.page.ready","jquery.makeCollapsible","mediawiki.toc","mediawiki.searchSuggest","ext.gadget.teahouse","ext.gadget.ReferenceTooltips","ext.gadget.watchlist-notice","ext.gadget.DRN-wizard","ext.gadget.charinsert","ext.gadget.refToolbar","ext.gadget.extra-toolbar-buttons","ext.gadget.switcher","ext.centralauth.centralautologin","mmv.head","mmv.bootstrap.autostart","ext.popups","ext.visualEditor.desktopArticleTarget.init","ext.visualEditor.targetLoader","ext.eventLogging","ext.wikimediaEvents","ext.navigationTiming","ext.uls.compactlinks","ext.uls.interface","ext.quicksurveys.init","ext.centralNotice.geoIP",
"ext.centralNotice.startUp","skins.vector.js"];</script>
<script>(RLQ=window.RLQ||[]).push(function(){mw.loader.implement("user.tokens@0tffind",function($,jQuery,require,module){/*@nomin*/mw.user.tokens.set({"editToken":"+\\","patrolToken":"+\\","watchToken":"+\\","csrfToken":"+\\"});
});});</script>
</html>
`
