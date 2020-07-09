# Sessions

This module provides simple, extensible session management infrastructure with built in support for memory and RDBMS backed sessions.

The module provides support for any ANSI SQL compliant database out of the box and it can be extended to support any other type of data stores such as NoSQL.

An important aspirational goal of this module is to provide comparable functionality with gorilla sessions with 75% less code.

The module will support following types of sessions:
- Standard: Typical sessions where session id is stored in cookie and corresponding session is stored on server. The data on server can be in stored in memory or any other persistent store such as databases.
- JWT: The session data is serialized as JWT and stored at client. The JWT can be stored in a cookie or in local storage at client. If stored local storage, JWT is sent to server as a header field in subsequent requests. 
- Hybrid JWT: The session id is stored as a JWT token on client side and it is used to lookup session data at server. This is similar to "Standard" option where session id is JWT token.