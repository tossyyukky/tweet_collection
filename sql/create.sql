CREATE TABLE tweets (
	id int unsigned auto_increment,
	content varchar(140) not null,
	username varchar(48) not null,
	tweeted timestamp not null,
	PRIMARY KEY(id)
);
