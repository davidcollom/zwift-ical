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


activate :data_source do |c|
    c.root  = "https://us-or-rly101.zwift.com/api/public/events"
    c.sources = [
        {
          alias: "events",
          path: "/upcoming",
          type: :json
        }
      ]
end

sports = @app.data.events.collect{|e| e['sport'].downcase }.uniq

sports.each do |sport|
    sport_events = @app.data.events.select{|e| e['sport']==sport.upcase}

    proxy "/#{sport}/rides", "ical", locals: {e: ZwiftCal.new(sport_events.select{|e| e['eventType'] == 'GROUP_RIDE'} ) }
    proxy "/#{sport}/workouts", "ical", locals: {e: ZwiftCal.new(sport_events.select{|e| e['eventType'] == 'GROUP_WORKOUT' } ) }
    proxy "/#{sport}/races", "ical", locals: {e: ZwiftCal.new(sport_events.select{|e| e['eventType'] == 'RACE' } ) }

    # Get all Unique Tags
    sport_events.collect{|e| e['tags']}.flatten.select{|t| !t.include?('=') }.sort.uniq do |t|
        proxy "/#{sport}/tag/#{t}", "ical", locals: {e: ZwiftCal.new(sport_events.select{|e| e['tags'].include?(t)} ) }
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
