# Sembot - Dokumentacja Biznesowa

Witaj w centrum dokumentacji procesów biznesowych i workflow technicznych Sembot.

## 🎯 Co znajdziesz w tej dokumentacji?

### 📚 Wytyczne i standardy
Poznaj zasady tworzenia dokumentacji biznesowej z naciskiem na diagramy Mermaid i czytelne procesy.

### 🔄 Procesy biznesowe  
Szczegółowe opisy procesów biznesowych z diagramami przepływu, rolami i odpowiedzialnościami.

### ⚙️ Workflow techniczne
Dokumentacja workflow deweloperskich, CI/CD, API i integracji systemowych.

### 🔗 Integracje
Opisy integracji z systemami zewnętrznymi, webhook'i i kolejki komunikatów.

## 🚀 Szybki start

1. **Nowy proces biznesowy?** → Sprawdź [wytyczne dokumentacji](/guides/)
2. **Tworzysz diagram?** → Zobacz [przykłady Mermaid](/examples/mermaid-examples)  
3. **Implementujesz feature?** → Przeczytaj [workflow techniczne](/workflow/)

## 📊 Przykład diagramu Mermaid

```mermaid
flowchart TD
    A[Użytkownik] --> B{Zalogowany?}
    B -->|Tak| C[Dashboard]
    B -->|Nie| D[Formularz logowania]
    D --> E[Weryfikacja danych]
    E -->|Poprawne| C
    E -->|Błędne| D
    C --> F[Funkcje systemu]
    
    classDef success fill:#d4edda,stroke:#155724
    classDef error fill:#f8d7da,stroke:#721c24  
    classDef process fill:#cce7ff,stroke:#0066cc
    
    class C,F success
    class D,E error
    class A,B process
```

## 🎨 Konwencje stylizacji

Wszystkie diagramy używają spójnej palety kolorów:
- 🟢 **Zielony** - procesy zakończone sukcesem
- 🔴 **Czerwony** - błędy i obsługa wyjątków  
- 🔵 **Niebieski** - standardowe kroki procesu
- 🟡 **Żółty** - procesy w trakcie/oczekujące

---

*Dokumentacja jest aktualizowana automatycznie przy każdym deploy na środowisko beta*