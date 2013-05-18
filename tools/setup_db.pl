#!/usr/bin/perl

use v5.10;
use strict;
use File::Basename;
use MongoDB;
use Digest;
use boolean;

my $dirname = dirname(__FILE__);

say "#"x100;
print <<EOF;
This script will create required collections for mongodb and will create admin user
EOF
say "#"x100;

say "Please enter admin login";
my $login = <>;
chomp $login;

say "Please enter admin passwd";
my $pass = <>;
chomp $pass;

# parse config
open FILE, "<$dirname/../conf/app.conf" or die $!;

my @lines = <FILE>;
close FILE;

my ($host, $port, $base);

for (@lines) {
  if ($_ =~ m/^mongodb.host\s*=\s*(.*)$/) {
    $host = $1;
  }
  if ($_ =~ m/^mongodb.port\s*=\s*(.*)$/) {
    $port = $1;
  }
  if ($_ =~ m/^mongodb.base\s*=\s*(.*)$/) {
    $base = $1;
  }
}

die "Wrong app.conf file" if !$base or !$port or !$host; 

# parse constants
open FILE, "<$dirname/../app/constants/const.go" or die $!;
@lines = <FILE>;
close FILE;

my $user_collection;
for (@lines) {
  if ($_ =~ m/^const UsersCollectionName\s+=\s*"(.*)"$/) {
    $user_collection = $1;
  }
}

die "can't get user collection from constants/const.go" if !$user_collection;

my $client = MongoDB::MongoClient->new(host => "mongodb://$host:$port");
my $db = $client->get_database($base);
my $users =  $db->get_collection($user_collection);
my $out = $users->find({"username" => $login});

if ($out->next) { 
  say "user $login already created. Delete it and create new one?";
  for (;;) { 
    my $out = <>;
    chomp $out;
    if (!$out =~ m/^(yes|no)$/) {
      say 'Please type "yes" or "no"'
    }
    if ($out eq "yes") {
      $users->remove({"username" => $login});
      last;
    };
    exit 0 if $out eq "no";
  }
};

my $sha512 = Digest->new("SHA-512");
$sha512->add($pass);

$users->insert({
    "username" => $login, 
    "password" => $sha512->hexdigest,
    "isadmin" => true
  });
