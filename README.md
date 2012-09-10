Axiom Data Services Gateway
===========================
This project is a "data services gateway" which allows rich-client applications such as JavaScript or Dart based
browser apps to use arbitrary oData, SQL, or other data services in an authenticated, secure, user-configurable manner.

Configuration Layout
--------------------
 * axiomds - the main executable
 * axiomds.conf - optional, the main configuration file, can include other files

Configuration JSON File Overview
--------------------------------

    {
        namespace: string, optional, defaults to blank (not in a namespace) [TODO]
        include: string list of files or dirs, optional [TODO]
        configure: [{
                'RDBMS Data Source' : {
                    "Name" : identity of data source, string, required, panic if duplicate
                    "Driver" : driver ID of database (e.g. 'mysql'), string, required
                    "Parameters" : connection URL, string, required, alias 'Connection Parameters'
                }

                // can combine/mix other directives here, too
            }
            {
                include: string list of files or dirs, optional [TODO]
                // can combine/mix other directives here, too
            }
            {
                'SQL DML' : {
                    "Name" : string, required, panic if duplicate
                    "Table Name" : string, required
                    "Command" : enum, optional, defaults to 'Auto'
                    "Columns" : string list, optional, defaults to all columns in table, alias: 'Allowed Columns'
                    "Primary Key Column" : string or string list, optional, defaults to PKs in table, alias: 'Primary Key Columns'
                    "Return on Success" : enum, optional, defaults to 'Primary Key Columns'
                    "Where Criteria" : parameters, optional, defaults to all records [TODO]
                    "Validate on Client" : list of JavaScript for each column, optional [TODO?]
                }

                // can combine/mix other directives here, too
            }
            {
                'SQL DQL' : {
                    "Name" : string, required, panic if duplicate
                    "DQL" : string, required [TODO: allow data source-specific SQL]
                    "Parameters" : value source spec list, optional
                    "Return on Success" : enum, optional, defaults to 'Multiple Rows as Hash List'
                }

                // can combine/mix other directives here, too
            }
        ]
    }