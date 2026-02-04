# Type-Glish

A terminal RPG where your English skills are your weapon. You fight monsters by typing combat actions, and a Dungeon Master (powered by an LLM) judges your grammar. Good grammar deals more damage. Bad grammar gets you hit.

Built for people who want to practice English in a less boring way (well, I'm trying not to be boring).

![Combat](images/Pasted%20image.png)

![Combat Result](images/Pasted%20image%20(2).png)

## Features

- Turn-based combat against various creatures
- Grammar scoring (1-10) determines damage dealt
- Enemy counter-attacks based on your score
- HP bars for you and enemies
- Path choice events between combats (typing heals you)
- Victory/defeat states
- Settings to switch LLM providers

## Quick Install

```bash
go build -o type-glish ./cmd/game
./type-glish
```

### Using Gemini

1. Get an API key from [Google AI Studio](https://aistudio.google.com/apikey)
2. Run the game and select "Settings" → "Use Gemini"
3. Enter your API key

The default model in the code is Gemini 3.0 Flash Preview. You can change it in your user config folder.

### Using llama.cpp

1. Install [llama.cpp](https://github.com/ggerganov/llama.cpp)
2. Start the server with a GGUF model:
   ```bash
   llama-server -hf [some GGUF model that your computer can handle]
   ```
   Add `-ngl 99` if you want to use your GPU.
   
3. Run the game and select "Use llama.cpp"

I haven't found the perfect CPU model yet. Suggestions welcome.

## Development

The game uses the [State pattern](https://refactoring.guru/design-patterns/state) which fits well with Bubbletea's ELM architecture (Model, View, Update).

### Adding a new screen

Create a new file in `internal/states/` implementing the `GameState` interface:

```go
type GameState interface {
    Init(ctx *game.Context) tea.Cmd
    Update(msg tea.Msg, ctx *game.Context) (GameState, tea.Cmd)
    View(ctx *game.Context) string
}
```

Each state handles its own input, renders its own view, and returns the next state on transitions.

### Project structure

```
internal/
├── game/       # Game context, player stats, enemy data
├── llm/        # LLM client, prompts, response types
├── states/     # All game states (combat, menu, etc.)
├── ui/         # Styles and UI helpers
└── tui/        # Main model that delegates to states
```

## License

MIT

