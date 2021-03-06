package pinboard

import (
	"encoding/json"
	"testing"
)

func TestPostsUnmarshallingBrokenDescription(t *testing.T) {
	brokenJSON := `{"href":"https:\/\/twitter.com\/LeishaRiddel\/status\/961609576641540097\/photo\/1","description":false,"extended":"Favorite tweet:\n\nM\u0360o\u032c\u033a\u0362t\u0317\u0331\u0333h\u031d\u0355\u032c\u0318\u0332\u032e\u0326e\u0355\u031e\u032b\u0332\u0330r\u0348\u0359\u031f\u031f\nc\u034f\u0349\u0349\u0354\u0333\u0325\u033a\u0330h\u031d\u0353\u0349\u0326\u0339\u035di\u031d\u0348\u0348\u0356\u0332\u031c\u0332\u035fl\u0347\u0355\u0323\u032c\u032d\u0354\u035fd\u0335\u0332r\u0320\u0330\u032d\u0329\u0331\u035a\u0316e\u0317\u031c\u032b\u033cn\nT\u033ah\u0489\u032d\u032d\u0330\u0325\u0359e\u0316\u032dy\u0326\u0353\u0330\u0347\u0356\u0361 \u1e77\u0335\u0332s\u033c\u0333\u0319\u0329\u035de \u0331m\u0318\u031f\u00fd e\u0329\u033a\u0324\u032a\u0323\u0345y\u033a\u0355\u031c\u031e\u0356\u0330\u032e\u0229\u0356\u0324\u0355\u033a\u0325\u032cs to\u0489 \u031d\u034d\u0326\u0356\u033a\u0345p\u0348\u031e\u0355\u031elay\u0330\u033a\u034d\u033a\u0319\u032c\u0333 \u034f\u032c\u0359\u033c\u0347w\u0334\u0319\u033b\u032ci\u0338\u0331\u0325\u034e\u0316\u0347t\u0316\u0332\u0316\u032d\u032b\u035dh\u0318\u0320 \u0329\u0331\u031d\u032b\u032a\u0353\u0323m\u0489\u0331\u031d\u0329\u034e\u031e\u0353\u0323e M\u0335\u0339\u033c\u0329\u0319\u0318\u032ao\u0318\u0362t\u033c\u1e2b\u034e\u0318e\u0319r I\u0319\u0348\u0347\u032c\u033a\u0323\u032e\u0358t \u0322\u032b\u0330\u031e\u0323h\u0335\u0349\u034d\u032c\u0333\u034e\u0333u\u0355r\u0318\u0323\u0324\u032b\u0324\u0345t\u0330\u0315s pic.twitter.com\/CKUDNx8Maz\n\n\u2014 Leisha Riddel (@LeishaRiddel) February 8, 2018","meta":"ce6199bac28b896012f067cb64fd7226","hash":"9ad3e9d77a12d21903f9ccafda592d84","time":"2018-02-08T20:03:06Z","shared":"no","toread":"no","tags":"IFTTT Twitter FAVES"}`
	var p post

	err := json.Unmarshal([]byte(brokenJSON), &p)
	if err != nil {
		t.Error(err)
	}
	if p.Description != "" {
		t.Error("Description is not empty string")
	}
	if p.Hash != "9ad3e9d77a12d21903f9ccafda592d84" {
		t.Error("Wrong hash")
	}
}

func TestPostsUnmarshallingOkDescription(t *testing.T) {
	okJSON := `{"href":"https:\/\/www.goodsfromjapan.com\/demekin-kingyo-sukui-goldfish-pack-p-1832.html","description":"Demekin Kingyo Sukui Goldfish Pack of 100 | Goods From Japan","extended":"","meta":"0becff1bafc8c4d347f19065ab44c2b4","hash":"984c059c4bdf506c05aa03ad2721a34d","time":"2018-03-31T07:26:25Z","shared":"no","toread":"yes","tags":"IFTTT Pocket"}`
	var p post

	err := json.Unmarshal([]byte(okJSON), &p)
	if err != nil {
		t.Error(err)
	}
	if p.Description != "Demekin Kingyo Sukui Goldfish Pack of 100 | Goods From Japan" {
		t.Error("Description is incorrect")
	}
	if p.Hash != "984c059c4bdf506c05aa03ad2721a34d" {
		t.Error("Wrong hash")
	}
}
