# ShakeSearch Challenge

Welcome to the Pulley Shakesearch Challenge! This repository contains a simple web app for searching text in the complete works of Shakespeare.

## Prerequisites

To run the tests, you need to have [Go](https://go.dev/doc/install) and [Docker](https://docs.docker.com/engine/install/) installed on your system.

## Your Task

Your task is to fix the underlying code to make the failing tests in the app pass. There are 3 frontend tests and 3 backend tests, with 2 of each currently failing. You should not modify the tests themselves, but rather improve the code to meet the test requirements. You can use the provided Dockerfile to run the tests or the app locally. The success criteria are to have all 6 tests passing.

## Instructions

<img width="404" alt="image" src="https://github.com/ProlificLabs/shakesearch/assets/98766735/9a5b96b5-0e44-42e1-8d6e-b7a9e08df9a1">

*** 

**Do not open a pull request or fork the repo**. Use these steps to create a hard copy.

1. Create a repository from this one using the "Use this template" button.
2. Fix the underlying code to make the tests pass
3. Include a short explanation of your changes in the readme or changelog file
4. Email us back with a link to your copy of the repo

## Running the App Locally


This command runs the app on your machine and will be available in browser at localhost:3001.

```bash
make run
```

## Running the Tests

This command runs backend and frontend tests.

Backend testing directly runs all Go tests.

Frontend testing run the app and mochajs tests inside docker, using internal port 3002.

```bash
make test
```

Good luck!

Changes:

1. To resolve the CORS issue when running locally in a browser, I implemented the enableCors function.
2. I enhanced the logic in Lookup to prevent "index out of range" errors when slicing the CompleteWorks string. This improvement addresses potential issues when idx-250 is less than 0 or idx+250 exceeds the length of CompleteWorks.
3. I noticed that closely situated search terms could result in duplicate text snippets being returned. To remedy this, I integrated logic to merge overlapping or adjacent segments in the Search method.
4. Upon encountering a failing Go test, specifically "expected 20 results for query," I inferred the necessity for pagination. Consequently, I introduced a pagination mechanism that initially returns 20 results and stores the remainder in a global variable. Additionally, I implemented a "Load More" button and a new /loadMore endpoint. This endpoint successively returns 20 results from the global variable, removing the returned data until it's depleted.
5. To address case sensitivity issues, I replaced Lookup with FindAllIndex, enabling the use of regular expressions for case-insensitive searches.
6. I created a Docker image to run test.js and made minor adjustments to the "Load More" button's ID, which further refined the pagination functionality established in step 4 and fixed Load More test case.
7. The test 'should return search results for "romeo, wherefore art thou"' was also corrected following the switch from Lookup to FindAllIndex.