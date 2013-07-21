// rss_api.go
//
// "THE PIZZA-WARE LICENSE" (derived from "THE BEER-WARE LICENCE"):
// <whoami@dev-urandom.eu> wrote these files. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you think
// this stuff is worth it, you can buy me a pizza in return.

package rss

import (
	"strings"
	"testing"
)

const TEST_FEED = "<rss version=\"2.0\"><channel><title>/int/ Minecraft News</title><description>This autism is to strong for me</description><link>https://dev-urandom.eu</link><lastBuildDate>20 Jul 13 00:31 +0200</lastBuildDate><generator></generator><item><title>Texasball is a cunt</title><link></link><description></description><author>greenkitten</author><category></category><pubDate>19 Jul 13 06:29 +0200</pubDate></item><item><title>Texasball has returned, prepare for autism levels exceeding safe limit</title><link></link><description></description><author>Texasball</author><category></category><pubDate>19 Jul 13 05:29 +0200</pubDate></item><item><title>Shakomatic got on today :O</title><link></link><description></description><author>I_can_haz_News</author><category></category><pubDate>18 Jul 13 22:19 +0200</pubDate></item><item><title>test</title><link></link><description></description><author>PPA</author><category></category><pubDate>18 Jul 13 17:18 +0200</pubDate></item><item><title>I love having sex with men.</title><link></link><description></description><author>Enton</author><category></category><pubDate>17 Jul 13 22:37 +0200</pubDate></item><item><title>pls unban</title><link></link><description></description><author>greenkitten</author><category></category><pubDate>13 Jul 13 13:18 +0200</pubDate></item><item><title>The Yevs at spawn are horrible, they need to be removed from premises. Support the crusade or join them in their grave!</title><link></link><description></description><author>areteee</author><category></category><pubDate>13 Jul 13 01:58 +0200</pubDate></item><item><title>The Tetown Games will award Darynu the title of champion at 8PM tonight if no other competitors show up by 759PM (or in ruffly 9 hours)US central time</title><link></link><description></description><author>McFuzzer</author><category></category><pubDate>08 Jul 13 18:12 +0200</pubDate></item><item><title>I&#39;m not gay, but there&#39;s something about taking a dick up the ass that really gets me going</title><link></link><description></description><author>Enton</author><category></category><pubDate>08 Jul 13 06:45 +0200</pubDate></item><item><title>test</title><link></link><description></description><author>AlphaBernd</author><category></category><pubDate>07 Jul 13 21:48 +0200</pubDate></item><item><title>tetown games sign ups will end on tomorrow. i or kitty4fun will be on after 5;pm us central time for any questions</title><link>https://dev-urandom.eu</link><description></description><author>McFuzzer</author><category></category><pubDate>Fri, 05 Jul 2013 20:16:31 +0200</pubDate></item><item><title>i die, horatio</title><link>https://dev-urandom.eu</link><description></description><author>v1adimirr</author><category></category><pubDate>Tue, 02 Jul 2013 06:51:47 +0200</pubDate></item><item><title>on  an unrelated note contact me if you know how to dispose a body</title><link>https://dev-urandom.eu</link><description></description><author>K_Chris</author><category></category><pubDate>Tue, 02 Jul 2013 06:51:39 +0200</pubDate></item><item><title>No wait, I just killed him, sorry</title><link>https://dev-urandom.eu</link><description></description><author>greenkitten</author><category></category><pubDate>Tue, 02 Jul 2013 06:50:21 +0200</pubDate></item><item><title>im not dead fuck all of you reading this</title><link>https://dev-urandom.eu</link><description></description><author>v1adimirr</author><category></category><pubDate>Tue, 02 Jul 2013 06:49:45 +0200</pubDate></item><item><title>enton is still gay</title><link>https://dev-urandom.eu</link><description></description><author>K_Chris</author><category></category><pubDate>Tue, 02 Jul 2013 06:49:06 +0200</pubDate></item><item><title>THE SIGN UP LIST IS UNDER BATTLES UNDER UP COMING EVENTS RULES WILL BE POSTED</title><link>https://dev-urandom.eu</link><description></description><author>kitty4fun</author><category></category><pubDate>Tue, 02 Jul 2013 01:09:57 +0200</pubDate></item><item><title>The Games at tetown are one weak away. i need to know how many participants by the 6th</title><link>https://dev-urandom.eu</link><description></description><author>kitty4fun</author><category></category><pubDate>Tue, 02 Jul 2013 01:08:06 +0200</pubDate></item><item><title>Enton is gay</title><link>https://dev-urandom.eu</link><description></description><author>Sikandar[IRC]</author><category></category><pubDate>Tue, 25 Jun 2013 18:26:00 +0200</pubDate></item><item><title>the games at tetown on the 8th have a sign up list for persons interested and persons who want to join</title><link>https://dev-urandom.eu</link><description></description><author>kitty4fun</author><category></category><pubDate>Tue, 25 Jun 2013 05:07:26 +0200</pubDate></item><item><title>Mod abuse is my favourite meme, I practice it daily!!</title><link>https://dev-urandom.eu</link><description></description><author>Entoron</author><category></category><pubDate>Sun, 23 Jun 2013 15:50:57 +0200</pubDate></item><item><title>competition at tetown will be on the 8th of next month</title><link>https://dev-urandom.eu</link><description></description><author>kitty4fun</author><category></category><pubDate>Sun, 23 Jun 2013 00:58:50 +0200</pubDate></item><item><title>authorize</title><link>https://dev-urandom.eu</link><description></description><author>Darynu</author><category></category><pubDate>Sat, 22 Jun 2013 10:10:52 +0200</pubDate></item><item><title>I hereby authorise my diamonds to be distributed to the people</title><link>https://dev-urandom.eu</link><description></description><author>caBst|IRC</author><category></category><pubDate>Tue, 18 Jun 2013 11:45:58 +0200</pubDate></item><item><title>Competition the 25th in Condeura at noon central time, and one at Mari Colosseum on the 8th of next month.</title><link>https://dev-urandom.eu</link><description></description><author>I_Can_Haz_News</author><category></category><pubDate>Sat, 15 Jun 2013 23:43:03 +0200</pubDate></item><item><title>Condeura is hosting a competition in 2 weeks see page for more info</title><link>https://dev-urandom.eu</link><description></description><author>I_Can_haz_News</author><category></category><pubDate>Tue, 11 Jun 2013 18:09:27 +0200</pubDate></item><item><title>i hav an erection</title><link>https://dev-urandom.eu</link><description></description><author>greenkitten</author><category></category><pubDate>Sun, 02 Jun 2013 23:27:29 +0200</pubDate></item><item><title>The Rhodesian Confederation recognizes the Battkhort-Breshikan Empire and K_Chris as its King-Emperor.</title><link>https://dev-urandom.eu</link><description></description><author>caBst|IRC</author><category></category><pubDate>Sun, 02 Jun 2013 23:26:37 +0200</pubDate></item><item><title>&gt;caBst doesn&#39;t play &gt;names a barely active player as new leader of Rhodes. Well sure, you&#39;re not to blame anymore now, HUEHAUHUEHUAHUEHUE</title><link>https://dev-urandom.eu</link><description></description><author>Dookie</author><category></category><pubDate>Fri, 31 May 2013 09:21:57 +0200</pubDate></item><item><title>Darynu has made a claim of the Breshikan throne</title><link>https://dev-urandom.eu</link><description></description><author>Darynu</author><category></category><pubDate>Thu, 30 May 2013 06:16:52 +0200</pubDate></item><item><title>will take roses as donations to whereisthe1pice</title><link>https://dev-urandom.eu</link><description></description><author>I_can_has_news</author><category></category><pubDate>Thu, 30 May 2013 06:15:50 +0200</pubDate></item><item><title>Rhodes wont have a new nation bid in a new map, if we ever get one. Dookie is to blame.</title><link>https://dev-urandom.eu</link><description></description><author>caBst|IRC</author><category></category><pubDate>Wed, 29 May 2013 23:15:42 +0200</pubDate></item><item><title>Badface has broken ground on its second city!</title><link>https://dev-urandom.eu</link><description></description><author>FlakeSe</author><category></category><pubDate>Wed, 29 May 2013 23:11:45 +0200</pubDate></item><item><title>Work on Nouvelle Reims is proceeding steadily! You&#39;re free to join the project!</title><link>https://dev-urandom.eu</link><description></description><author>Ranshiin</author><category></category><pubDate>Wed, 29 May 2013 22:36:23 +0200</pubDate></item><item><title>i proclame my self ruler of condeura</title><link>https://dev-urandom.eu</link><description></description><author>Minuke_00</author><category></category><pubDate>Sat, 25 May 2013 06:24:40 +0200</pubDate></item><item><title>Rosenmann didn&#39;t confess about stealing intellectual property</title><link>https://dev-urandom.eu</link><description></description><author>mazznoff2</author><category></category><pubDate>Fri, 24 May 2013 18:58:23 +0200</pubDate></item><item><title>Please contact Sikandar[IRC] if interested in helping build a recreational server V2 mod (representing some history and nations of the server)</title><link>https://dev-urandom.eu</link><description></description><author>Sikandar[IRC]</author><category></category><pubDate>Fri, 24 May 2013 08:28:10 +0200</pubDate></item><item><title>ghoul isn&#39;t dead</title><link>https://dev-urandom.eu</link><description></description><author>K_Chris</author><category></category><pubDate>Thu, 23 May 2013 19:06:39 +0200</pubDate></item><item><title>kitty4fun has left the rule of Condeura to its Barrons to decide</title><link>https://dev-urandom.eu</link><description></description><author>kitty4fun</author><category></category><pubDate>Thu, 23 May 2013 15:37:39 +0200</pubDate></item><item><title>Shako has left the server; Breshik merged with Bracave; P_P_A Count of Nova Avence and Bydlograd</title><link>https://dev-urandom.eu</link><description></description><author>PPA</author><category></category><pubDate>Sun, 19 May 2013 00:40:45 +0200</pubDate></item><item><title>United Nations declares the True Rhodesian Confederacy a illegitimate state</title><link>https://dev-urandom.eu</link><description></description><author>caBst|IRC</author><category></category><pubDate>Thu, 16 May 2013 22:23:47 +0200</pubDate></item><item><title>caBst|IRC declared traitor by the True Rhodesian Confederacy</title><link>https://dev-urandom.eu</link><description></description><author>Uxbridge</author><category></category><pubDate>Thu, 16 May 2013 22:22:17 +0200</pubDate></item><item><title>The Rhodesian Confederation has declared war on Breshik</title><link>https://dev-urandom.eu</link><description></description><author>caBst|IRC</author><category></category><pubDate>Thu, 16 May 2013 22:08:02 +0200</pubDate></item><item><title>Also the rest of the KCO</title><link>https://dev-urandom.eu</link><description></description><author>Sikandar[IRC]</author><category></category><pubDate>Thu, 16 May 2013 22:06:13 +0200</pubDate></item><item><title>The Khanate has declared war on Breshik.</title><link>https://dev-urandom.eu</link><description></description><author>greenkitten</author><category></category><pubDate>Thu, 16 May 2013 22:05:34 +0200</pubDate></item><item><title>v1adimirr is still alive</title><link>https://dev-urandom.eu</link><description></description><author>AlphaBernd</author><category></category><pubDate>Tue, 14 May 2013 03:01:05 +0200</pubDate></item><item><title>Kitty4fun is attempting a bank in Condeura</title><link>https://dev-urandom.eu</link><description></description><author>kitty4fun</author><category></category><pubDate>Mon, 13 May 2013 23:06:11 +0200</pubDate></item><item><title>&lt;@Enton&gt; i&#39;m not the one suffering from unwarranted self importance</title><link>https://dev-urandom.eu</link><description></description><author>areteee</author><category></category><pubDate>Sun, 12 May 2013 21:06:47 +0200</pubDate></item><item><title>&lt;caBst|IRC&gt;Sorry but I&#39;m not needy of attention.</title><link>https://dev-urandom.eu</link><description></description><author>mazznoff</author><category></category><pubDate>Sat, 11 May 2013 16:33:56 +0200</pubDate></item><item><title>I got gf</title><link>https://dev-urandom.eu</link><description></description><author>Dookie</author><category></category><pubDate>Sat, 11 May 2013 10:41:54 +0200</pubDate></item><item><title>Enton is the best way to deliver anything... As long as its minecraft things</title><link>https://dev-urandom.eu</link><description></description><author>Eren1</author><category></category><pubDate>Fri, 10 May 2013 18:16:55 +0200</pubDate></item><item><title>People of Condeura see wiki page changes</title><link>https://dev-urandom.eu</link><description></description><author>kitty4fun</author><category></category><pubDate>Fri, 10 May 2013 17:45:52 +0200</pubDate></item><item><title>If you want an edit made to the family tree please suggest it on nations:familytree NOT ANYWHERE ELSE. Thank you t:genealogist</title><link>https://dev-urandom.eu</link><description></description><author>greenkitten</author><category></category><pubDate>Fri, 10 May 2013 13:41:24 +0200</pubDate></item><item><title>Everything is fine, cute and adorable in Kriegstein, which is great</title><link>https://dev-urandom.eu</link><description></description><author>Dookie</author><category></category><pubDate>Thu, 09 May 2013 21:51:30 +0200</pubDate></item><item><title>the Empire of Battkhortostan has declared war on Kovanje, renounced Koer autonomy</title><link>https://dev-urandom.eu</link><description></description><author>Chris</author><category></category><pubDate>Thu, 09 May 2013 21:45:49 +0200</pubDate></item><item><title>Holy Banana Empire now claiming West West Battx, including the rest of Battx</title><link>https://dev-urandom.eu</link><description></description><author>Enton</author><category></category><pubDate>Thu, 09 May 2013 19:06:05 +0200</pubDate></item><item><title>Battkhortostan has announced the full annexation of Breshik</title><link>https://dev-urandom.eu</link><description></description><author>Sikandar[IRC]</author><category></category><pubDate>Thu, 09 May 2013 19:05:28 +0200</pubDate></item><item><title>Rhodesian Confederation stands with RCF Member New Walschor in case of embassy breach during Batt.x occupation</title><link>https://dev-urandom.eu</link><description></description><author>caBst|IRC</author><category></category><pubDate>Thu, 09 May 2013 18:46:39 +0200</pubDate></item><item><title>Koinu pedal</title><link>https://dev-urandom.eu</link><description></description><author>Koinu</author><category></category><pubDate></pubDate></item><item><title>Factories are popping up on server!</title><link>https://dev-urandom.eu</link><description></description><author>kitty4fun</author><category></category><pubDate></pubDate></item><item><title>Al-Iskandariya is under siege by a coalition of Kovanje Boers and West West Battkhortistani anarchists.</title><link>https://dev-urandom.eu</link><description></description><author>greenkitten</author><category></category><pubDate></pubDate></item><item><title>that feel when still no gf</title><link>https://dev-urandom.eu</link><description></description><author>ijs</author><category></category><pubDate></pubDate></item><item><title>still no gf</title><link>https://dev-urandom.eu</link><description></description><author>areteee</author><category></category><pubDate></pubDate></item><item><title>caBst|IRC was last seen 8 hours 45 minutes 22 seconds ago: Yeah because it sucks</title><link>https://dev-urandom.eu</link><description></description><author>AlphaBernd</author><category></category><pubDate></pubDate></item><item><title>HAL now has a news system that can be connected to the wiki</title><link>https://dev-urandom.eu</link><description></description><author>AlphaBernd</author><category></category><pubDate></pubDate></item></channel></rss>"
const TEST_FILE = "test.rss"

func TestRSSParser(t *testing.T) {
	sr := strings.NewReader(TEST_FEED)

	rss, err := ParseFromReader(sr)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if rss.Version != "2.0" {
		t.Fatalf("Parsed wrong Version: %s", rss.Version)
	}

	if len(rss.Channel) == 0 || len(rss.Channel) > 1 {
		t.Fatalf("Wrong Channel length: %d", len(rss.Channel))
	}

	if len(rss.Channel[0].Item) == 0 {
		t.Fatal("Couldn't parse items")
	}

	if rss.Channel[0].Item[0].Author != "greenkitten" {
		t.Fatal("Did not parse first author correctly")
	}
}

func TestAddItem(t *testing.T) {
	sr := strings.NewReader(TEST_FEED)

	rss, err := ParseFromReader(sr)
	if err != nil {
		t.Fatalf("%v", err)
	}

	rss.AddItem("TestTitle", "TestLink", "TestDescription", "TestAuthor", "TestCategory")

	if rss.Channel[0].Item[0].Title != "TestTitle" {
		t.Fatal("Did not add Item correctly")
	}
}

func TestWriteFile(t *testing.T) {
	sr := strings.NewReader(TEST_FEED)

	rss, err := ParseFromReader(sr)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := rss.WriteToFile(TEST_FILE); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestFileParser(t *testing.T) {
	rss, err := ParseFromFile(TEST_FILE)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if rss.Version != "2.0" {
		t.Fatalf("Parsed wrong Version: %s", rss.Version)
	}

	if len(rss.Channel) == 0 || len(rss.Channel) > 1 {
		t.Fatalf("Wrong Channel length: %d", len(rss.Channel))
	}

	if len(rss.Channel[0].Item) == 0 {
		t.Fatal("Couldn't parse items")
	}

	if rss.Channel[0].Item[0].Author != "greenkitten" {
		t.Fatal("Did not parse first author correctly")
	}
}
