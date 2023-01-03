require_relative 'zwiftcal'

activate :dotenv

# Activate gzip compression
activate :gzip

ignore 'vendor'
ignore '.env'
ignore 'ical'

set :css_dir, 'css'
set :js_dir, 'js'
set :images_dir, 'img'
set :fonts_dir,  "fonts"

set :file_watcher_ignore,[
    /^bin(\/|$)/,
    /^\.bundle(\/|$)/,
    /^.c9(\/|$)/,
    /^node_modules(\/|$)/,
    /^\.sass-cache(\/|$)/,
    /^\.cache(\/|$)/,
    /^\.git(\/|$)/,
    /^\.gitignore$/,
    /\.DS_Store/,
    /^\.rbenv-.*$/,
    /^Gemfile$/,
    /^Gemfile\.lock$/,
    /~$/,
    /(^|\/)\.?#/,
    /^tmp\//
  ]

# binding.pry

main_events = ZwiftCal.events()
sports = main_events.collect{|e| e['sport'].downcase }.sort.uniq

sports.each do |sport|
    sport_events = main_events.select{|e| e['sport']==sport.upcase}

    proxy "/#{sport}/rides", "ical", locals: {e: ZwiftCal.new(sport_events.select{|e| e['eventType'] == 'GROUP_RIDE'} ) }
    proxy "/#{sport}/workouts", "ical", locals: {e: ZwiftCal.new(sport_events.select{|e| e['eventType'] == 'GROUP_WORKOUT' } ) }
    race_events = sport_events.select{|e| e['eventType'] == 'RACE' }
    proxy "/#{sport}/races", "ical", locals: {e: ZwiftCal.new(race_events) }
    # [A..E].do |klass|
    #   proxy "/#{sport}/races/#{klass}", "ical", locals: {e: ZwiftCal.new(race_events.select{|e| e['eventSubgroups'].include? 'RACE' } ) }
    # end

    # Get all Unique Tags
    sport_events.collect{|e| e['tags']}.flatten.select{|t| !t.include?('=') }.sort.uniq.each do |tag|
        tag_events = (ZwiftCal.events(tag: tag)||[]).select{|e| e['sport'] == sport.upcase}
        puts "#{sport.upcase}: Found #{tag_events.length} Events for #{tag}"
        proxy "/#{sport}/tag/#{tag}", "ical", locals: {e: ZwiftCal.new(tag_events) }
    end
end



configure :server do
end

configure :development do
end

configure :production do
end

# Build-specific configuration
configure :build do
end
