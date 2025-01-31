package util

/*type testCase struct {
	user string
	text string
}

func TestAnonymizer(t *testing.T) {
	msgs := []slack.Message{
		{
			Msg: slack.Msg{
				Username: "rvanderp",
				Text:     "connect to: test.test.com",
			},
		},
		{
			Msg: slack.Msg{
				Username: "rvanderp2",
				Text:     "connect to: 192.168.1.2",
			},
		},
		{
			Msg: slack.Msg{
				Username: "rvanderp3",
				Text:     "mac address is: c6:56:20:11:2d:a5",
			},
		},
		{
			Msg: slack.Msg{
				Username: "rvanderp4",
				Text:     "mac address is: c6:56:20:11:2d:a5 <@rvanderp4>",
			},
		},
	}

	expectedText := []testCase{
		{
			user: "op",
			text: "xxxxxxxxxxxxxxxxxxxxxxxxx",
		},
		{
			user: "contributor_2",
			text: "connect to: x-ipv4-0000000001-x",
		},
		{
			user: "contributor_3",
			text: "mac address is: x-mac-0000000001-x",
		},
		{
			user: "contributor_4",
			text: "mac address is: x-mac-0000000001-x <@contributor_4>",
		},
	}

	msgs = AnonymizeMessages(msgs)

	for i, msg := range msgs {
		if msg.Text != expectedText[i].text {
			t.Errorf("expected %s, got %s", expectedText[i].text, msg.Text)
		}
		if msg.Username != expectedText[i].user {
			t.Errorf("expected %s, got %s", expectedText[i].user, msg.Username)
		}
	}
}*/
