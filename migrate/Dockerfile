FROM ruby:3

WORKDIR /app/migrate

COPY ./migrate/Gemfile /app/migrate
COPY ./migrate/Gemfile.lock /app/migrate
RUN bundle install

ENTRYPOINT ["rake"]
CMD ["db:migrate"]