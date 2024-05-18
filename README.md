# Purpose
Each week my manager wants a detailed written summary of work that is more fine grained than the standup board.  My memory often fails me by the end of the week, causing unnecessary stress and wasted time. This tool allows me to efficiently record my daily tasks and generate a comprehensive summary using ChatGPT, making it easy to compile and submit my weekly reports.

# Stack
- 🚀 Go
- ☕ BubbleTea withh Bubbles
- 🤖 ChatGPT Completions API
- 🗃️ SQLite

# Running
1) Create a file called config.json in the /config folder, and fill it in with your ChatGPT API credentials. And example can be found in config.example.json
2) Create a file called prompt.txt in the /config folder, and add whatever prompt you want to feed to ChatGPT.  You can start by using the one found in example-prompt.txt
3) `go run .`

# Why the ChatAPI Completetions endpoint?
This is a [legacy endpoint](https://platform.openai.com/docs/api-reference/completions), but I chose to use it simply because my company has our own dedicated instance of ChatGPT for internal use, which uses the Completions endpoint.  This meant I could prototype without needing to pay for API costs.  If you want to use a different endpoint the code can be found in the llm package.