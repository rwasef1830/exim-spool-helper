package main

import (
	"strings"
	"testing"
)

func TestGetProcessedSpoolHeaderContent(t *testing.T) {
	input := `
tamer.fouad@customer.com

tamer.fouad@customer.com

265P Received: from [1.2.3.4] (port=55804 helo=AmirSabry)
        by server.com with esmtpsa  (TLS1.2) tls TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
        (Exim 4.96)
        (envelope-from <contracts@customer.com>)
        id 1pgPvo-0004lc-34;
        Sun, 26 Mar 2023 14:57:17 +0200
130T To: "Person yehia" <person.yehia@customer.com>,
        "'Person Fouad'" <person.fouad@customer.com>,
        <'Person.yousry@customer.com'>
094C Cc: "bamya person" <b.p@customer.com>,
        "elnagar person" <elnagar.p@customer.com>
065  Subject: Some subject
038  Date: Sun, 26 Mar 2023 14:57:07 +0200
061I Message-ID: <d95fe2$797ca600$6c75f2@customer.com>
018  MIME-Version: 1.0
085  Content-Type: multipart/mixed;
        boundary="----=_NextPart_000_0085_01D95FF3.3D057600"
033  X-Mailer: Microsoft Outlook 16.0
047  Thread-Index: Adlf4lsKCp5idQovQKHg9oiIdw==
024  Content-Language: en-us
031F From: contracts@customer.com




`

	expected := `
tamer.fouad@customer.com
        
243P Received: from [1.2.3.4] (port=55804 helo=AmirSabry)
	by server.com with esmtpsa  (TLS1.2) tls TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
	(Exim 4.96)
	(envelope-from <contracts@customer.com>)
	id 1pgPvo-0004lc-34;
	Sun, 26 Mar 2023 14:57:17 +0200
127T To: "Person yehia" <person.yehia@customer.com>,
	"'Person Fouad'" <person.fouad@customer.com>,
	<'Person.yousry@customer.com'>
082C Cc: "bamya person" <b.p@customer.com>,
	"elnagar person" <elnagar.p@customer.com>
022  Subject: Some subject
038  Date: Sun, 26 Mar 2023 14:57:07 +0200
050I Message-ID: <d95fe2$797ca600$6c75f2@customer.com>
018  MIME-Version: 1.0
085  Content-Type: multipart/mixed;
	boundary="----=_NextPart_000_0085_01D95FF3.3D057600"
033  X-Mailer: Microsoft Outlook 16.0
043  Thread-Index: Adlf4lsKCp5idQovQKHg9oiIdw==
024  Content-Language: en-us
029F From: contracts@customer.com
`

	content, err := getProcessedSpoolHeaderContent("test", strings.NewReader(input))
	if err != nil {
		t.Errorf("Test failed: %s", err)
	}

	if content != expected {
		t.Failed()
	}
}
