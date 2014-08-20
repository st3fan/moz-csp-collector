-- create-tables.sql

create table Site (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  Hostname varchar unique not null
);

create table UserAgent (
  Id serial not null unique primary key,
  Value varchar unique not null
);

create table DocumentUri (
  Id serial not null unique primary key,
  Value varchar unique not null
);

create table Referrer (
  Id serial not null unique primary key,
  Value varchar unique not null
);

create table BlockedUri (
  Id serial not null unique primary key,
  Value varchar unique not null
);

create table ViolatedDirective (
  Id serial not null unique primary key,
  Value varchar unique not null
);

create table OriginalPolicy (
  Id serial not null unique primary key,
  Value varchar unique not null
);

create table ScriptSource (
  Id serial not null unique primary key,
  Value varchar unique not null
);

create table ScriptSample (
  Id serial not null unique primary key,
  Value varchar unique not null
);

create table Report (
  Id serial not null unique primary key,
  Created timestamp without time zone default (now() at time zone 'utc'),
  SiteId integer references Site (Id),
  UserAgentId integer references UserAgent (Id),
  RemoteIp varchar not null,
  DocumentUriId integer not null references DocumentUri (Id),
  ReferrerId integer references Referrer (Id),
  BlockedUriId integer references BlockedUri (Id),
  ViolatedDirectiveId integer not null references ViolatedDirective (Id),
  OriginalPolicyId integer references OriginalPolicy (Id),
  ScriptSourceId integer references ScriptSource (Id),
  ScriptSampleId integer references ScriptSample (Id)
);
