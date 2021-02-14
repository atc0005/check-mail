// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/check-mail
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

var testStrings = []struct {
	original string
	modified string
}{
	{
		"What's new from Gravity Forms? 🚀 September Edition",
		"What's new from Gravity Forms? ✂ September Edition",
	},
	{
		"Celebrate WooCommerce Day with 40% off everything 🎉",
		"Celebrate WooCommerce Day with 40% off everything ✂",
	},
	{
		"SALE: 30% off WooCommerce.com until midnight 👻",
		"SALE: 30% off WooCommerce.com until midnight ✂",
	},
	{
		"🔔 Last chance to take 35% off WooCommerce.com",
		"✂ Last chance to take 35% off WooCommerce.com",
	},
	{
		"⏰ Just a few more hours! ⏰",
		"⏰ Just a few more hours! ⏰",
	},
	{
		"Early Cyber Deals are Here! 🐱🎉🐶",
		"Early Cyber Deals are Here! ✂✂✂",
	},
	{
		// within the supported character utf8mb3 character set (surprisingly)
		"⚡ LAST CHANCE to Join ⭐ Give the Gift of Learning",
		"⚡ LAST CHANCE to Join ⭐ Give the Gift of Learning",
	},
	{
		"🎁 Final Reminder: 50% off Go Courses ends tomorrow!",
		"✂ Final Reminder: 50% off Go Courses ends tomorrow!",
	},
	{
		"🤪 CYBER MONDAY MADNESS",
		"✂ CYBER MONDAY MADNESS",
	},
	{
		"It's Cyber Monday 👾 Take 55% OFF everything today only!",
		"It's Cyber Monday ✂ Take 55% OFF everything today only!",
	},
	{
		"🚨🍩CYBER MONDAY ONLY🚨🍩 FREE Original Glazed Dozen offer inside",
		"✂✂CYBER MONDAY ONLY✂✂ FREE Original Glazed Dozen offer inside",
	},
	{
		"📢 Deals have been extended though 12/4!",
		"✂ Deals have been extended though 12/4!",
	},
	{
		"How do we know the Earth is round? 🌎",
		"How do we know the Earth is round? ✂",
	},
	{
		"30% Sitewide Cyber Monday Savings 😱 Only Once a Year!",
		"30% Sitewide Cyber Monday Savings ✂ Only Once a Year!",
	},
	{
		"❓What do you want to learn about AWS?",
		"❓What do you want to learn about AWS?",
	},
	{
		// leading *3* Unicode code points: U+1F385, U+1F3FD, U+1F384
		"🎅🏽🎄 TONS of Christmas Printables + LAST Day to Join TLC & MMB",
		"✂✂✂ TONS of Christmas Printables + LAST Day to Join TLC & MMB",
	},
	{
		"💻 Cyber Monday! Free shipping today only and deals throughout the week!",
		"✂ Cyber Monday! Free shipping today only and deals throughout the week!",
	},
	{
		"Adam, cozy up to soup season. 🥣",
		"Adam, cozy up to soup season. ✂",
	},
	{
		"Don’t Miss Your Chance to Save! ⏰🎁💰",
		"Don’t Miss Your Chance to Save! ⏰✂✂",
	},
	{
		"🎄 More Christmas Printables + 🕎 Hanukkah Pack",
		"✂ More Christmas Printables + ✂ Hanukkah Pack",
	},
	{
		"Don’t miss the WooCommerce.com marketplace sale! ⏰",
		"Don’t miss the WooCommerce.com marketplace sale! ⏰",
	},
	{
		"Adam, Get New Year Ready! 🎆🍾🎉",
		"Adam, Get New Year Ready! ✂✂✂",
	},
	{
		"🎉 WooCommerce 3.2 is here, bringing improved coupons, extension management, and more!",
		"✂ WooCommerce 3.2 is here, bringing improved coupons, extension management, and more!",
	},
	{
		// contains non-breaking spaces before the emoji (U+00A0, \u00a0)
		"Brand new eBooks and videos only $10 each! Get 'em while they're hot. 🔥",
		"Brand new eBooks and videos only $10 each! Get 'em while they're hot. ✂",
	},
	{
		// contains non-breaking spaces before the emoji (U+00A0, \u00a0)
		"Only 72 hours left! Don't miss out on $10 titles & $25 bundles! ⌚",
		"Only 72 hours left! Don't miss out on $10 titles & $25 bundles! ⌚",
	},
	{
		// contains non-breaking spaces before the emoji (U+00A0, \u00a0)
		"Find the perfect bundle-of-3 for $25 📚",
		"Find the perfect bundle-of-3 for $25 ✂",
	},
	{
		// contains non-breaking spaces before the emoji (U+00A0, \u00a0)
		"Win a golden ticket to WooConf in Seattle  😍",
		"Win a golden ticket to WooConf in Seattle  ✂",
	},
	{
		// seems to be acceptable to MySQL/MariaDB with utf8 (utf8mb3)
		"ő",
		"ő",
	},
	{
		// seems to be acceptable to MySQL/MariaDB with utf8 (utf8mb3)
		"szükséges információk megküldése",
		"szükséges információk megküldése",
	},
	{
		// "smiley of doom"
		// https://emojipedia.org/grinning-face-with-smiling-eyes/
		"\xF0\x9F\x98\x81",
		"✂",
	},
	{
		"\xF4\x80\x80\x80",
		"✂",
	},
	{
		"🆑",
		"✂",
	},
	{
		"🐳",
		"✂",
	},
	{
		"\xF0\xA0\xAE\x9F",
		"✂",
	},
	{
		"\xE2\x9C\xA8",
		"\xE2\x9C\xA8",
	},
	{
		"✂",
		"✂",
	},
	{
		"\xE2\x9C\x82",
		"\xE2\x9C\x82",
	},
	{
		"☀",
		"☀",
	},
	{
		"🍓🍆Fruit & Veggies Beginning Letter Sounds + 🐧Penguin Playdough Mat!",
		"✂✂Fruit & Veggies Beginning Letter Sounds + ✂Penguin Playdough Mat!",
	},
	{
		"Sip some stressed spelled backwards 🍰 🍵 🍫",
		"Sip some stressed spelled backwards ✂ ✂ ✂",
	},
}
