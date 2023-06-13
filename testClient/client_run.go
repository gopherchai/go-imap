package main_1

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/textproto"
	"os"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
)

var (
	client *imapclient.Client
)

func setup(t *log.Logger) {
	c, err := imapclient.DialTLS("imap.qq.com:993", nil)
	if err != nil {
		t.Fatalf("dial error:%+v", err)
		return
	}
	client = c

}

func main2() {
	t := log.New(os.Stderr, "", log.Lshortfile)
	setup(t)
	err := client.Login("codexyz@foxmail.com", "nwmcmqrqioyabjei").Wait()
	t.Printf("%+v\n", err)
	cmd := client.List("", "%", nil)
	// &imap.ListOptions{
	// 	SelectSubscribed:     true,
	// 	SelectRemote:         true,
	// 	SelectRecursiveMatch: true,
	// 	ReturnSubscribed:     true,
	// 	ReturnChildren:       true,

	// 	 ReturnStatus: &imap.StatusOptions{
	// 	 	NumMessages: true,
	// 	 	NumUnseen:   true,
	// 	 },
	// }

	list, err := cmd.Collect()
	if err != nil {
		t.Println(err)
	}
	for _, item := range list {
		mbox := item
		t.Printf("%+v\n", mbox)
		//	t.Printf("Mailbox %q contains %v messages (%v unseen)", mbox.Mailbox, mbox.Status.NumMessages, mbox.Status.NumUnseen)

	}
	t.Println("error", cmd.Close())
	selectedMbox, err := client.Select("其他文件夹", nil).Wait()
	t.Println("error", err)
	t.Printf("%+v", selectedMbox)
	selectedMbox, err = client.Select("gptbox", &imap.SelectOptions{
		ReadOnly: false,
	}).Wait()
	t.Println("error", err)
	t.Printf("%+v", selectedMbox)

	idleCmd, err := client.Idle()
	if err != nil {
		t.Println("error", err)
		return
	}

	err = idleCmd.Close()
	err2 := idleCmd.Wait()
	if err != nil || err2 != nil {
		t.Println("error", err, err2)
		return
	}

	seqSet := imap.SeqSetRange(1, 3)
	fetchItems := []imap.FetchItem{
		imap.FetchItemBody,
		&imap.FetchItemBinarySectionSize{},
		//imap.FetchItemRFC822,
		//internal.FetchItemRFC822,
		imap.FetchItemEnvelope,
		imap.FetchItemBodyStructure,
		&imap.FetchItemBodySection{},
		&imap.FetchItemBinarySection{},
		// imap.FetchItemBodyStructure,
		imap.FetchItemUID,
		imap.FetchItemFlags,
		imap.FetchItemRFC822Size,
	}

	// client.Search(&imap.SearchCriteria{
	// 	NotFlag: []imap.Flag{imap.FlagSeen},
	// }, nil)
	messages, err := client.Fetch(seqSet, fetchItems, nil).Collect()
	if err != nil {
		t.Fatalf("failed to fetch first message in gptbox: %v", err)
	}
	//t.Printf("subject of first message in INBOX: %+v\n", messages[0])

	for _, msg := range messages {
		t.Printf("%+v,%+v\n", msg.Envelope, msg.UID)
	}
	msg := messages[0]

	//msg.BodyStructure
	v := msg.BodyStructure
	t.Printf("%+v\n", v)

	mp, ok := v.(*imap.BodyStructureMultiPart)

	if ok {
		t.Printf("%+v\n", mp)
		t.Printf("%+v\n", mp.Extended)
		for _, v := range mp.Children {
			v.MediaType()
			t.Printf("%+v,\n", v)
			p, ok := v.(*imap.BodyStructureSinglePart)
			if ok {
				t.Printf("%+v\n", p)
				t.Printf("%+v\n", p.Text)

			}
		}
	}

	for _, v := range msg.BinarySection {

		t.Printf("bdata %v\n", string(v))
	}

	for k, v := range msg.BodySection {
		t.Printf("%+v\n", k.HeaderFields)
		bdata, _ := json.Marshal(k)
		t.Printf("bdata %v\n", string(bdata))
		t.Printf("bdata %v\n", string(v))
		mr, _ := mail.CreateReader(bytes.NewBuffer(v))
		ct, params, err := mr.Header.ContentType()
		t.Printf("%+v,%+v,%+v\n", ct, params, err)
		message.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
			return nil, nil
		}
		part, err := mr.NextPart()
		t.Printf("%+v,%+v,%+v\n", part.Header.Get("Content-Transfer-Encoding"), part.Header.Get("Content-Type"))
		bodyBytes := make([]byte, msg.RFC822Size)
		part.Body.Read(bodyBytes)
		t.Printf("%+v\n", string(bodyBytes))
		buf := bytes.NewBuffer(v)
		r := textproto.NewReader(bufio.NewReader(buf))

		header, err := r.ReadMIMEHeader() //Content-Type:[multipart/alternative; boundary="----=_NextPart_6460B1E8_C21D8B40_06C29985"]
		t.Printf("%+v,%+v\n", header, err)
		//buf2 := bytes.NewBuffer(v)
		//textproto.NewMultipartReader()
		//mmp := multipart.NewReader(buf2, "----=_NextPart_6460B1E8_C21D8B40_06C29985")
		//part, err := mmp.NextPart()

		if strings.Contains(part.Header.Get("Content-Type"), "text/plain") {

		} else if strings.Contains(part.Header.Get("Content-Type"), "text/html") {
			part.Header.Get("Content-Transfer-Encoding")
			//map[Content-Transfer-Encoding:[base64] Content-Type:[text/plain; charset="utf-8"]]
			//根据不同的字符集进行解码，不同的content-type进行读取

		}
		t.Printf("%+v,%+v,%+v,%+v\n", part, err, part.Header, msg.BinarySectionSize)
		data := make([]byte, msg.RFC822Size)

		t.Printf("%+v,%+v\n", string(data), msg.RFC822Size)

	}
	bdata, _ := json.Marshal(msg.BodySection)

	t.Printf("bdata %v\n", string(bdata))
	seqSet = imap.SeqSetNum(1)
	fetchItems = []imap.FetchItem{
		imap.FetchItemBody,
		//imap.FetchItemRFC822,

		imap.FetchItemEnvelope,
		// imap.FetchItemBodyStructure,
		// imap.FetchItemUID,
		// imap.FetchItemFlags,
		// imap.FetchItemRFC822Size,
	}
	msgs, err := client.UIDFetch(seqSet, fetchItems, &imap.FetchOptions{}).Collect()
	t.Printf("%+v,%+v", msgs[0], err)
	msg = msgs[0]

	t.Printf("%+v\n", msg.Envelope)

	//msg.BodyStructure
	v = msg.BodyStructure
	t.Printf("%+v\n", v)

	mp, ok = v.(*imap.BodyStructureMultiPart)

	if ok {
		t.Printf("%+v\n", mp)

		for _, v := range mp.Children {

			t.Printf("%+v\n", v)
			p, ok := v.(*imap.BodyStructureSinglePart)
			if ok {
				t.Printf("%+v\n", p)
				t.Printf("%+v\n", p.Text)

			}
		}
	}

	for _, v := range msg.BinarySection {

		t.Printf("bdata %v\n", string(v))
	}

	for k, v := range msg.BodySection {
		t.Printf("%+v\n", k.HeaderFields)
		bdata, _ := json.Marshal(k)
		t.Printf("bdata %v\n", bdata)
		t.Printf("bdata %v\n", string(v))
	}
	bdata, _ = json.Marshal(msg.BodySection)

	t.Printf("bdata %v\n", string(bdata))
	sres, err := client.Search(&imap.SearchCriteria{

		NotFlag: []imap.Flag{imap.FlagSeen},
	}, nil).Wait()
	if err != nil {
		t.Panic(err)
	}

	messages, _ = client.Fetch(sres.All, fetchItems, nil).Collect()
	for _, msg := range messages {
		t.Printf("%+v,%+v\n", msg.Envelope, msg.UID)
	}

	err = client.Store(imap.SeqSetNum(messages[0].SeqNum), &imap.StoreFlags{
		Op:     imap.StoreFlagsAdd,
		Silent: true,
		Flags:  []imap.Flag{imap.FlagSeen},
	}, &imap.StoreOptions{}).Wait()
	t.Printf("%+v,\n", err)

	//每秒扫一次邮箱，将未读邮件读取出来
}
