---
category: api
title: Dokumentacja API
author: System Documentation
description: Dokumentacja REST API, endpoints, autentykacja i przykłady użycia
---

# 🌐 Dokumentacja API

Ten katalog zawiera dokumentację REST API systemu Sembot.

## 📋 Co powinno się tutaj znaleźć?

### 1. Dokumentacja API
- **API Reference** - Pełna dokumentacja endpoints
- **Authentication** - Autentykacja i autoryzacja
- **Rate Limiting** - Limity i throttling
- **Versioning** - Wersjonowanie API
- **Error Handling** - Obsługa błędów

### 2. Grupy endpoints
- **Campaigns API** - Zarządzanie kampaniami
- **Products API** - Zarządzanie produktami
- **Workflows API** - Wykonywanie workflow
- **Tasks API** - Zarządzanie zadaniami
- **Connections API** - Połączenia zewnętrzne
- **Reports API** - Raporty i analytics

### 3. Przykłady użycia
- Request/Response examples
- cURL commands
- SDK examples (PHP, JavaScript)

## 🎯 Struktura plików

```
api/
├── README.md                           (ten plik)
├── DOC_API_Overview.md                 - Przegląd API
├── DOC_API_Authentication.md           - Autentykacja
├── DOC_API_Rate_Limiting.md            - Rate limiting
├── DOC_API_Errors.md                   - Obsługa błędów
├── endpoints/
│   ├── DOC_Campaigns_API.md            - Campaigns endpoints
│   ├── DOC_Products_API.md             - Products endpoints
│   ├── DOC_Workflows_API.md            - Workflows endpoints
│   ├── DOC_Tasks_API.md                - Tasks endpoints
│   ├── DOC_Connections_API.md          - Connections endpoints
│   └── DOC_Reports_API.md              - Reports endpoints
└── examples/
    ├── DOC_API_Examples_cURL.md        - cURL examples
    ├── DOC_API_Examples_PHP.md         - PHP SDK examples
    └── DOC_API_Examples_JavaScript.md  - JS examples
```

## 📝 Format dokumentu endpoint

```markdown
---
category: api
title: Nazwa API - Endpoints
author: Twoje Imię
description: Dokumentacja endpoints dla [zasób]
---

## Przegląd
[Opis grupy endpoints]

## Authentication
[Wymagane uprawnienia]

## Endpoints

### GET /api/v1/resource
Pobiera listę zasobów.

**Request:**
```http
GET /api/v1/resource?page=1&limit=50
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 50
  }
}
```

**Errors:**
- `401 Unauthorized` - Brak lub nieprawidłowy token
- `403 Forbidden` - Brak uprawnień
- `429 Too Many Requests` - Rate limit exceeded

### POST /api/v1/resource
Tworzy nowy zasób.

[...podobna struktura...]
```

## 🎯 Konwencje

### HTTP Methods
- `GET` - Pobieranie danych
- `POST` - Tworzenie zasobu
- `PUT/PATCH` - Aktualizacja zasobu
- `DELETE` - Usuwanie zasobu

### Response Codes
- `200 OK` - Sukces
- `201 Created` - Zasób utworzony
- `204 No Content` - Sukces bez treści
- `400 Bad Request` - Błąd walidacji
- `401 Unauthorized` - Brak autentykacji
- `403 Forbidden` - Brak uprawnień
- `404 Not Found` - Zasób nie znaleziony
- `422 Unprocessable Entity` - Błąd walidacji
- `429 Too Many Requests` - Rate limit
- `500 Internal Server Error` - Błąd serwera

### Pagination
```json
{
  "data": [...],
  "links": {
    "first": "url",
    "last": "url",
    "prev": null,
    "next": "url"
  },
  "meta": {
    "current_page": 1,
    "from": 1,
    "to": 50,
    "total": 100,
    "per_page": 50,
    "last_page": 2
  }
}
```

## 📝 Jak dodać dokumentację?

1. Utwórz plik `DOC_*_API.md`:
```yaml
---
category: api
title: Nazwa API
author: Twoje Imię
description: Opis
---
```

2. Dokumentuj wszystkie endpoints z:
   - Request format
   - Response format
   - Error codes
   - Przykłady

3. Build docs:
```bash
cd .doc && ./build-docs.sh
```

## 🔗 Powiązane

- `/domain/` - Business logic za API
- `/guides/` - Przewodniki użycia API
- `/integration/` - Integracje z API
