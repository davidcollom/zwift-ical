require 'sinatra'
require_relative 'zwiftcal'

if ENV['REDIS_HOST'] != ''
  require 'lightly'
  $cache = Lightly.new life: '14h'
else
  require 'redis-store'
  $cache = Redis::Store.new host: ENV['REDIS_HOST'], port: ENV['REDIS_PORT'] || 6379, db: ENV['REDIS_DB'] || 0
end

class App < ::Sinatra::Base

  set :root, '/app/'

  set :sessions, false
  set :static, true
  set :public_folder, '/app/public/'
  set :logging, true

  get '/' do
    erb :index
  end

  ## Cycling Routes
  get '/cycling/search' do
    content_type 'text/calendar'

    error 403 do
      "Must provide a q paramater"
    end if params[:q].empty?
    $cache.disable
    events = $cache.get("cycling:events:search") do
      puts "Downloading fresh data from zwift..."
      ZwiftCal.events
    end
    $cache.enable
    halt(404) unless events.any?
    ZwiftCal.new( events.select{|e| e['name'].include? params[:q]} ).to_ical
  end

  get '/cycling/rides' do
    content_type 'text/calendar'

    events = $cache.get("cycling:events:groups") do
      puts "Downloading fresh data from zwift..."
      ZwiftCal.events.select{|e| e['eventType'] == 'GROUP_RIDE'}
    end
    halt(404) unless events.any?
    ZwiftCal.new(events).to_ical
  end

  get '/cycling/workouts' do
    content_type 'text/calendar'
    events = $cache.get("cycling:events:workouts") do
      puts "Downloading fresh data from zwift..."
      ZwiftCal.events.select{|e| e['eventType'] == 'GROUP_WORKOUT' && e['sport']=='CYCLING'}
    end
    halt(404) unless events.any?
    ZwiftCal.new(events).to_ical
  end

  get '/cycling/races' do
    content_type 'text/calendar'

    events = $cache.get("cycling:events:races") do
      puts "Downloading fresh data from zwift..."
      ZwiftCal.events.select{|e| e['eventType'] == 'RACE' && e['sport']=='CYCLING'}
    end
    halt(404) unless events.any?
    ZwiftCal.new(events).to_ical
  end

  get '/cycling/tag/:tag(.:format)?' do
    content_type 'text/calendar'

    events = $cache.get("cycling:events:tags:#{params['tag']}") do
      puts "Downloading fresh data..."
      ZwiftCal.events(tag: params[:tag]).select{|e| e['sport']=='CYCLING'}
    end
    halt(404) unless events.any?
    ZwiftCal.new(events).to_ical
  end

  ## Running Routes
  get '/running/workouts' do
    content_type 'text/calendar'
    events = $cache.get("running:events:workouts") do
      puts "Downloading fresh data from zwift..."
      ZwiftCal.events.select{|e| e['eventType'] == 'GROUP_WORKOUT' && e['sport']=='RUNNING'}
    end
    halt(404) unless events.any?
    ZwiftCal.new(events).to_ical
  end

  get '/running/races' do
    content_type 'text/calendar'

    events = $cache.get("running:events:races") do
      puts "Downloading fresh data from zwift..."
      ZwiftCal.events.select{|e| e['eventType'] == 'RACE' && e['sport']=='RUNNING'}
    end
    halt(404) unless events.any?
    ZwiftCal.new(events).to_ical
  end

  get '/running/tag/:tag(.:format)?' do
    content_type 'text/calendar'

    events = $cache.get("running:events:tags:#{params['tag']}") do
      puts "Downloading fresh data..."
      ZwiftCal.events(tag: params[:tag]).select{|e| e['sport']=='RUNNING'}
    end
    halt(404) unless events.any?
    ZwiftCal.new(events).to_ical
  end

  get '/version' do
    ENV['VERSION'] || 'Unknown'
  end


  after do
    $cache.prune
  end

end
