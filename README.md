# BabelDB

_“The library will endure; it is the universe. As for us, everything has not been written; we are not turning into phantoms. We walk the corridors, searching the shelves and rearranging them, looking for lines of meaning amid leagues of cacophony and incoherence, reading the history of the past and our future, collecting our thoughts and collecting the thoughts of others, and every so often glimpsing mirrors, in which we may recognize creatures of the information.”_ ― Jorge Luis Borges, The Library of Babel

⚠️ BabelDB is an ongoing experimentation project.

**BabelDB** is an in-memory Website Database. BabelDB combines a programmatic data extraction engine with scheduling and data clustering. It offers a standard and lightweight SQL syntax and a powerful DSL for querying, searching and information retrieval.
**BabelDB** continuously ingests data from any pre-defined web source and allows you to query data with standard SQL. Also it provides its own query language: BabelQL, built on top of the engine to provide search capabilities such as full-text search, term and phrase matching, regex and more.

## Features

- [ ] Data collection scheduling
- [ ] Data clustering
- [ ] Links tagging
- [ ] Incrementally updated materialized views
- [ ] Pattern matching
- [ ] Deep collection
- [ ] Stream data into pre-defined sinks
- [ ] Define custom data collectors
- [ ] Semantic subscription
- [ ] Data discovery


## Motivation

From Wikipedia: 
```
...a database is an organized collection of data stored and accessed electronically...
```

#### Can Internet as a whole be considered a Database itself?

The internet is a vast space of information. Most of the information is free (which does not mean true) and accessible through browsers and search engines and dedicated tooling. Crawler & Scrapper bots are popular ways for automated data collection and indexing. Crawling is essentially what search engines do while scraping is an automated way of extracting specific datasets.
But when it comes to address a more specific use cases or non-technical users, sometimes this is not enough. 

For example: 
- I want to collect all news articles automatically and compare climate change narrative between site X and Y.
- I want to know how site X looked like 24 hours ago and retrieve only the updates.
- I want to keep track of companies that are environmentally friendly or have sustainability programs.
- I want to discover linked web resources which match with some pattern.
- I want to subscribe and be aware when certain semantic shows up in site X.

Ok!, technically speaking this is not too complex with the tooling we have access nowadays. But let's say I want a Marketing analyst with knowledge of SQL can do it.  

**BabelDB** is the experimental attempt to solve that! 😀


