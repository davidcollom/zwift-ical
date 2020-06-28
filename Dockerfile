FROM ruby:2.3.7
ARG VERSION=develop

RUN gem install bundler

COPY Gemfile /app/
WORKDIR /app
RUN bundle install

COPY . /app/

ENV PORT=3000 APP_ENV=production RACK_ENV=production VERSION=$VERSION
EXPOSE 3000

CMD ["puma", "config.ru","-p","3000"]
