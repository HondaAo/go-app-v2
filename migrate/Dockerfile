FROM ruby:3

WORKDIR /app/migrate

COPY ./migrate/Gemfile /app/migrate
COPY ./migrate/Gemfile.lock /app/migrate
RUN bundle install

RUN useradd migrate
USER migrate

ENTRYPOINT ["rake"]
CMD ["db:migrate"]