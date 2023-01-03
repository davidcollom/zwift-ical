require 'httparty'
require 'retries'
require 'zache'
require 'icalendar'
require 'icalendar/tzinfo'

class ZwiftCal
  include ::HTTParty

  base_uri 'https://us-or-rly101.zwift.com/api/public/events'
  format :json
  headers 'Content-Type' => 'application/json'
  logger ::Logger.new STDOUT, :debug, :curl


  WORLDS = [
    'Watopia',
    'Richmond',
    'London',
    'New York',
    'Innsbuck',
    'Bologna TT',
    'Yorkshire',
    'Crit City',
    'France',
    'Paris',
    'Makuri Islands'
  ].freeze

  TIMEZONE = "UTC" # Nasty hack to force UTC

  @@cache = Zache.new()

  class << self
    def events(limit: 200, tag: "", options: {})
      build_query = {tags: tag, limit: limit}
      options.merge!( {query: build_query} )
      @@cache.get(options.to_s) do
        puts "Fetching Events...#{options}"
        with_retries(:max_tries => 10, :base_sleep_seconds => 0.1, :max_sleep_seconds => 2.0) do
          self.get("/upcoming", options )
        end
      end
    end
  end

  def initialize(events)
    @cal = Icalendar::Calendar.new
    @cal.prodid = "Zwift Calendar - by David Collom"
    @events = events
    parse_events
  end

  def parse_events
    @events.each do |event|
      @cal.event do |e|
        e.uid         = event['id'].to_s
        e.summary     = event_summary(event)
        e.description = event['description']
        e.dtstart     = Icalendar::Values::DateTime.new DateTime.parse( event['eventStart'] ), 'tzid' => TIMEZONE
        e.dtend       = Icalendar::Values::DateTime.new calculate_end(event), 'tzid' => TIMEZONE
        e.url         = event['imageUrl'] if event['imageUrl']!=''
        e.ip_class    = "PUBLIC"
        e.last_modified =  Icalendar::Values::DateTime.new DateTime.parse( Time.now.utc.to_s ), 'tzid' => TIMEZONE
        # e.append_attach   Icalendar::Values::Uri.new("ftp://host.com/novo-procs/felizano.exe", "fmttype" => "application/binary")
      end
    end
  end

  def to_s
    @cal.to_ical
  end

  def to_ical
    @cal.to_ical
  end

  private
  def event_summary(event)
    "[#{world_name(event)}] #{event['name']}"
  end

  def calculate_end(event)
    # If event is duration or distance related
    if event['durationInSeconds'] == 0
      # puts "Assuming #{event['name']} is 1 hour long [#{event['durationInSeconds']}]"
      DateTime.parse( event['eventStart'] ).to_time + 3600 # Assume ~1 hour
    else
      # Add duration to event start
      # puts "#{event['name']} is #{event['durationInSeconds']} seconds long"
      (DateTime.parse( event['eventStart'] ).to_time) + event['durationInSeconds']
    end
  end

  def world_name(event)
    map_id = event['mapId'] || 0
    # puts "#{event['name']} is in #{event['mapId']}"
    if !WORLDS[ map_id-1 ].nil?
      WORLDS[map_id-1 ].to_s
    else
      "Unknown"
    end
  end

end
