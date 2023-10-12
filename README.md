# Snippetbox

## Summary

This repository contains my learnings from Alex Edwards' brilliant [Let's Go](https://lets-go.alexedwards.net/) book. While I have stuck quite closely with the project defined in the book, I have also made deviations and additions to: (1) further challenge myself, and (2) explore more of the popular packages used by the greater Go community.

Please consider this repository as a work-in-progress — there are more things I want to add and make changes to.

## Deviations/Additions

1. I have added a `Makefile` adapted from a [blog post]((https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects)) by Alex Edwards to make running mundane commands quicker. Notably, the `Makefile` includes a command to run the project with live reloading which automatically reruns the development server whenever a code change is detected.
2. I have added a `package.json` file that includes [Prettier](https://prettier.io/) and the [prettier-plugin-go-template](https://github.com/NiklasPor/prettier-plugin-go-template) plugin for it that formats Go HTML templates. I find it quite difficult to maintain the templates otherwise due to the Go templating syntax.
3. I use the [GoDotEnv](https://github.com/joho/godotenv) package for reading `.env` files which I use for providing default configuration values, rather than hard-coding them.
4. I use the [Chi](https://github.com/go-chi/chi) library for routing instead of [httprouter](https://github.com/julienschmidt/httprouter). Chi offers many benefits over httprouter like clean middleware chaining and route grouping while remaining stdlib compliant.
5. I use [pgx](https://github.com/jackc/pgx) rather than the standard [database/sql](https://pkg.go.dev/database/sql). This is because I decided to use Postgres instead of MySQL (which Alex uses in the book). Additionally, pgx features performance benefits due to its focus on Postgres rather than SQL databases in general.
6. I have added metrics via [Prometheus](https://github.com/prometheus/client_golang) to get the codebase closer to a production-ready state. I visualise the metrics using [Grafana Cloud](https://grafana.com/products/cloud/).