package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestGetProcessedSpoolHeaderContent1(t *testing.T) {
	input := `
1phRtv-0001yu-2y-H
mailnull 47 12
-received_time_usec .921469
-received_time_complete 1680081327.932839
-helo_name smtp14-fra-sp2.mta.salesforce.com
-host_address [85.222.158.221]:9433
-host_name smtp14-fra-sp2.mta.salesforce.com
-interface_address [54.36.238.13]:25
-received_protocol esmtps
-aclm 0 1
1
-aclm 1 8
iabgroup
-body_linecount 112
-max_received_linelength 250
-frozen 1680081334
-spam_bar ++
-spam_score 2.6
-spam_score_int 26
-tls_cipher TLS1.2:ECDHE-RSA-AES256-GCM-SHA384:256
-tls_ourcert -----BEGIN CERTIFICATE-----\nMIIF6DCCBNCgAwIBAgIRAOIkmlGG3NqF/PqK/zKMdkYwDQYJKoZIhvcNAQELBQAw\ncjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAlRYMRAwDgYDVQQHEwdIb3VzdG9uMRUw\nEwYDVQQKEwxjUGFu$
-tls_resumption B
-tls_ver TLS1.2
XX
1
tamer.fouad@customer.com

garbage

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

	expected := `1phRtv-0001yu-2y-H
mailnull 47 12
<contracts@customer.com>
-received_time_usec .921469
-received_time_complete 1680081327.932839
-helo_name smtp14-fra-sp2.mta.salesforce.com
-host_address [85.222.158.221]:9433
-host_name smtp14-fra-sp2.mta.salesforce.com
-interface_address [54.36.238.13]:25
-received_protocol esmtps
-aclm 0 1
1
-aclm 1 8
iabgroup
-body_linecount 112
-max_received_linelength 250
-frozen 1680081334
-spam_bar ++
-spam_score 2.6
-spam_score_int 26
-tls_cipher TLS1.2:ECDHE-RSA-AES256-GCM-SHA384:256
-tls_ourcert -----BEGIN CERTIFICATE-----\nMIIF6DCCBNCgAwIBAgIRAOIkmlGG3NqF/PqK/zKMdkYwDQYJKoZIhvcNAQELBQAw\ncjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAlRYMRAwDgYDVQQHEwdIb3VzdG9uMRUw\nEwYDVQQKEwxjUGFu$
-tls_resumption B
-tls_ver TLS1.2
XX
1
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
029F From: contracts@customer.com`

	content, err := getProcessedSpoolHeaderContent("test", strings.NewReader(input))
	if err != nil {
		t.Errorf("Test failed: %s", err)
	}

	if diff := cmp.Diff(content, expected); diff != "" {
		fmt.Println(diff)
		t.Fail()
	}
}

func TestGetProcessedSpoolHeaderContent2(t *testing.T) {
	input := `
1phRtv-0001yu-2y-H
mailnull 47 12
-received_time_usec .921469
-received_time_complete 1680081327.932839
-helo_name smtp14-fra-sp2.mta.salesforce.com
-host_address [85.222.158.221]:9433
-host_name smtp14-fra-sp2.mta.salesforce.com
-interface_address [54.36.238.13]:25
-received_protocol esmtps
-aclm 0 1
1
-aclm 1 8
iabgroup
-body_linecount 112
-max_received_linelength 250
-frozen 1680081334
-spam_bar ++
-spam_score 2.6
-spam_score_int 26
-tls_cipher TLS1.2:ECDHE-RSA-AES256-GCM-SHA384:256
-tls_ourcert -----BEGIN CERTIFICATE-----\nMIIF6DCCBNCgAwIBAgIRAOIkmlGG3NqF/PqK/zKMdkYwDQYJKoZIhvcNAQELBQAw\ncjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAlRYMRAwDgYDVQQHEwdIb3VzdG9uMRUw\nEwYDVQQKEwxjUGFu$
-tls_resumption B
-tls_ver TLS1.2
XX
1
tamer.fouad@customer.com

garbage

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
031F From: Test <contracts@customer.com>




`

	expected := `1phRtv-0001yu-2y-H
mailnull 47 12
<contracts@customer.com>
-received_time_usec .921469
-received_time_complete 1680081327.932839
-helo_name smtp14-fra-sp2.mta.salesforce.com
-host_address [85.222.158.221]:9433
-host_name smtp14-fra-sp2.mta.salesforce.com
-interface_address [54.36.238.13]:25
-received_protocol esmtps
-aclm 0 1
1
-aclm 1 8
iabgroup
-body_linecount 112
-max_received_linelength 250
-frozen 1680081334
-spam_bar ++
-spam_score 2.6
-spam_score_int 26
-tls_cipher TLS1.2:ECDHE-RSA-AES256-GCM-SHA384:256
-tls_ourcert -----BEGIN CERTIFICATE-----\nMIIF6DCCBNCgAwIBAgIRAOIkmlGG3NqF/PqK/zKMdkYwDQYJKoZIhvcNAQELBQAw\ncjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAlRYMRAwDgYDVQQHEwdIb3VzdG9uMRUw\nEwYDVQQKEwxjUGFu$
-tls_resumption B
-tls_ver TLS1.2
XX
1
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
036F From: Test <contracts@customer.com>`

	content, err := getProcessedSpoolHeaderContent("test", strings.NewReader(input))
	if err != nil {
		t.Errorf("Test failed: %s", err)
	}

	if diff := cmp.Diff(content, expected); diff != "" {
		fmt.Println(diff)
		t.Fail()
	}
}
