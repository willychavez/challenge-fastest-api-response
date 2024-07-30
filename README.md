# challenge-fastest-api-response

- **cmd/challenge-fastest-api-response/**: Main application entry point.
- **internal/api/**: Contains the logic for API calls.

## APIs Used

1. **BrasilAPI**: A Brazilian API that provides data about postal codes, among other things.
2. **ViaCEP**: Another popular API for Brazilian postal codes.

## How It Works

The program sends simultaneous requests to both APIs for a given CEP. It waits for the first response and displays the result, ignoring the slower response. The maximum waiting time is 1 second; if neither API responds within this period, a timeout error is reported.

## Running the Project

To run the project, ensure you have Go installed on your system. Then, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/willychavez/challenge-fastest-api-response.git
   ```
2. Navigate to the project directory:
   ```bash
   cd challenge-fastest-api-response
   ```
3. Run the project:
   ```bash
   go run ./cmd/challenge-fastest-api-response
   ```

## Example Output

```bash
Result from BrasilAPI: {CEP:01153000 State:SP City:São Paulo Neighborhood:Barra Funda Street:Rua Vitorino Carmilo}
```
Or in case of an error:

```bash
Error ViaCEP: Get "http://viacep.com.br/ws/01153000/json/": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

## Customizing the CEP

To test with a different CEP, modify the `cep` variable in the main function:

```go
cep := "01001000" // Replace with the desired CEP
```

## Dependencies

This project uses the standard Go library, including the `net/http` and `encoding/json` packages for making HTTP requests and parsing JSON responses.
