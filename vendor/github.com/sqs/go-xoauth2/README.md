go-xoauth2
====================

Go library for generating strings for use in the XOAUTH2 authentication scheme,
implemented by Google's Gmail IMAP and SMTP servers for OAuth2 authentication.


Usage
--------------------

    import (
           "xoauth2"
    )
    
    ...
    
    // returns unencoded string
    xoauth2.OAuth2String("alice", "some-access-token")
  
    // returns base64-encoded XOAUTH2 string, suitable for direct use in SASL XOAUTH2
    xoauth2.XOAuth2String("alice", "some-access-token")


More information
--------------------

See the XOAUTH2 protocol documentation for more information:
    https://developers.google.com/google-apps/gmail/xoauth2_protocol#the_sasl_xoauth2_mechanism
