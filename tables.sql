-- Sites that are allowed to report. Contains just the hostname.

create table Sites (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  Hostname varchar unique not null
);

create table Documents (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  Uri varchar unique not null
);

create table Referrers (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  Uri varchar unique not null
);

create table Blockers (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  Uri varchar unique not null
);

create table UserAgents (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  UserAgent varchar not null
);

create table Policies (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  Policy varchar not null
);

create table Directives (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  Directive varchar unique not null
);

create table Reports (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  -- Site this happened on
  SiteId integer references Sites (Id),
  -- Client info
  ClientIp varchar not null,
  ClientUserAgent integer references UserAgents (Id),
  -- This is the report
  DocumentId integer not null references Documents (Id),
  ReferrerId integer references Referrers (Id),
  BlockerId integer references Blockers (Id),
  ViolatedDirectiveId integer not null references Directives (Id),
  OriginalPolicyId integer references Policies (Id),
  -- Firefox specific
  ScriptSource varchar,
  ScriptSample varchar
);


CREATE OR REPLACE FUNCTION InsertDocument(DocumentUri VARCHAR) RETURNS INTEGER AS $$
BEGIN
  RETURN QUERY
  INSERT INTO Documents (Uri)
    SELECT DocumentUri
      WHERE NOT EXISTS (SELECT Id FROM Documents WHERE Uri = DocumentUri)
  RETURNING Id;
END; $$
LANGUAGE PLPGSQL;


CREATE OR REPLACE FUNCTION inc(val integer) RETURNS integer AS $$
BEGIN
RETURN val + 1;
END; $$
LANGUAGE PLPGSQL;


INSERT INTO Documents (Uri) SELECT DocumentUri WHERE NOT EXISTS (SELECT Id FROM Documents WHERE Uri = DocumentUri) RETURNING Id;
