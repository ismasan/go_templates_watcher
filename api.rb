require 'rubygems'
require 'sinatra'

put '/templates/:name' do |name|
  p [:recd, name, params]
end