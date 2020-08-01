package handlers

import (
	"testing"
)

func TestFilterStrings(t *testing.T) {
	inputs := []string{
		// Street names with context
		"This is a test, and replaced should be this: [[\"Jl. SMA Aek Kota Batu\",\"id\"],[\"Sumatera Utara\",\"de\"]], yes that is what should be replaced.",
		// Street names
		"[[\"а/д Вятка\",\"ru\"]]",
		// Shop object
		"[[[null,[\"5170389350216449417\",\"8960489397473601205\"]],[[[0.3817225694656372,0.411925345659256,3.330686807632446]]],[\"GENTS Maastricht\",\"nl\"],[\"Herrenmodengeschäft\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_2x.5.png\",[null,null,null,null,null,null,null,null,\"J2NkhgQ4yws\"]],[[null,[\"5170389356187074925\",\"938526169172417203\"]],[[[0.7049263119697571,0.4769878685474396,39.86749649047852]]],[\"Marks\",\"nl\"],[\"Eiscafé\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/restaurant_2x.5.png\",[null,null,null,null,null,null,null,null,\"YQfV2MfERP4\"]],[[null,[\"5170389355496820819\",\"1932532312540701841\"]],[[[0.8073962926864624,0.4982556402683258,11.67663192749023]]],[\"Café Falstaff\",\"nl\"],[\"Café\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/cafe_2x.5.png\",[null,null,null,null,null,null,null,null,\"x1CHoBLhAsM\"]],[[null,[\"5170389356238105859\",\"17180988966320210211\"]],[[[0.4690370857715607,0.4416649341583252,33.43334579467773]]],[\"MENners clothing Maastricht\",\"nl\"],[\"Herrenmodengeschäft\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_2x.5.png\",[null,null,null,null,null,null,null,null,\"W1UZ3AGjsQg\"]],[[null,[\"5170389038332097889\",\"6920768074393208644\"]],[[[0.459946870803833,0.4442543387413025,20.50216484069824]]],[\"MR \u0026 MRS Maastricht\",\"en\"],[\"Bekleidungsgeschäft\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_2x.5.png\",[null,null,null,null,null,null,null,null,\"2QXmksl57Bg\"]],[[null,[\"5170389356944885023\",\"12214692630577502547\"]],[[[0.9193480610847473,0.5108497142791748,7.628006458282471]]],[\"Dille \u0026 Kamille\",\"en\"],[\"Haushaltswarengeschäft\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_2x.5.png\",[null,null,null,null,null,null,null,null,\"1XJdkl7rBhk\"]],[[null,[\"5170389355515317381\",\"6033344794370014695\"]],[[[0.6844578981399536,0.4723083972930908,26.90214729309082]]],[\"Kinsjasa Schoenen\",\"nl\"],[\"Schuhgeschäft\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_2x.5.png\",[null,null,null,null,null,null,null,null,\"RZrMFcx3dK4\"]],[[null,[\"5170389355510483265\",\"6598220819907671338\"]],[[[0.7621901631355286,0.4912925064563751,18.60411262512207]]],[\"café b.for\",\"en\"],[\"Café\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/cafe_2x.5.png\",[null,null,null,null,null,null,null,null,\"tgFh0HWJGeY\"]],[[null,[\"5170389355498751757\",\"12648706420154025534\"]],[[[0.7725918889045715,0.492845892906189,16.37007904052734]]],[\"Lincherie\",\"en\"],[\"Dessousgeschäft\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_2x.5.png\",[null,null,null,null,null,null,null,null,\"zgpN4jk4h2w\"]],[[null,[\"5170389356220794673\",\"4798407148106230527\"]],[[[0.5379243493080139,0.4389545023441315,15.93551826477051]]],[\"Clio Jewelry\",\"en\"],[\"Juwelier\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_2x.5.png\",[null,null,null,null,null,null,null,null,\"XuWhwqnuu-E\"]],[[null,[\"5170389356190094439\",\"6415132636374069217\"]],[[[0.6754223704338074,0.4693296551704407,24.99129676818848]]],[\"Freddy Pant Room Maastricht\",\"sr\"],[\"Damenmodengeschäft\",\"de\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_2x.5.png\",[null,null,null,null,null,null,null,null,\"4v1I-bHYXGk\"]]],null",
	}
	outputs := []string{
		"This is a test, and replaced should be this: [[\"\",\"\"],[\"\",\"\"]], yes that is what should be replaced.",
		"[[\"\",\"\"]]",
		"[],null",
	}

	for i := range inputs {
		out := string(filterStrings([]byte(inputs[i])))
		if out != outputs[i] {
			t.Fatal("Expected\n", outputs[i], "\nbut got\n", out)
		}
	}
}
