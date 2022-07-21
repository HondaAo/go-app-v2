FROM ruby:3

WORKDIR /app/migrate

COPY ./migrate/Gemfile /app/migrate
COPY ./migrate/Gemfile.lock /app/migrate

COPY ./migrate  /app/migrate

RUN gem install bundler
RUN bundle install

ENTRYPOINT ["rake"]
CMD ["db:migrate"]