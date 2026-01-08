# Go Backend → Java Comparison Guide 📚

Understanding the Go backend by comparing it to Java concepts you already know.

---

## 1. Project Structure

### Go Structure
```
movie-booking/
├── cmd/main.go              # Entry point (like main() method)
├── api/v1/                  # API layer (like Controllers)
├── core/services/           # Business logic (like Services)
├── datastore/               # Data access (like Repositories)
└── core/model/              # Models (like Entities)
```

### Java Equivalent
```
movie-booking/
├── src/main/java/
│   ├── Application.java     # @SpringBootApplication (main class)
│   ├── controller/          # @RestController (API layer)
│   ├── service/             # @Service (business logic)
│   ├── repository/          # @Repository (data access)
│   └── model/               # @Entity (JPA entities)
```

**Key Difference:**
- **Go**: Explicit separation, manual dependency injection
- **Java**: Spring Framework handles DI with annotations

---

## 2. Entry Point: `cmd/main.go`

### Go Code
```go
func main() {
    apiFlag := flag.Bool("api", false, "Start API server")
    flag.Parse()
    
    if *apiFlag {
        startAPIServer(db)
    }
}
```

### Java Equivalent
```java
@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

**Concept:**
- Both are entry points
- Go uses command-line flags (`--api`)
- Java uses Spring Boot auto-configuration

---

## 3. Models/Entities

### Go Model (`core/model/models.go`)
```go
type Movie struct {
    ID          uint   `gorm:"column:id"`
    Title       string `gorm:"column:title"`
    Description string `gorm:"column:description"`
    Duration    int    `gorm:"column:duration_mins"`
}
```

### Java Entity
```java
@Entity
@Table(name = "movies")
public class Movie {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(name = "title")
    private String title;
    
    @Column(name = "description")
    private String description;
    
    @Column(name = "duration_mins")
    private Integer duration;
    
    // getters/setters
}
```

**Key Differences:**
- **Go**: Struct tags (`gorm:"column:title"`) for ORM mapping
- **Java**: Annotations (`@Column`, `@Entity`) for JPA mapping
- **Go**: No getters/setters (direct field access)
- **Java**: Encapsulation with private fields + getters/setters

---

## 4. Controllers (API Layer)

### Go Controller (`api/v1/controllers/controller.go`)
```go
func (c *Controller) LoginHandler(w http.ResponseWriter, r *http.Request) (*types.GenericAPIResponse, error) {
    // Parse request
    req, err := helpers.ValidateAndParseLoginRequest(r)
    if err != nil {
        return nil, errors.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    
    // Call service
    result, err := c.authService.Login(ctx, req.Email, req.Password)
    if err != nil {
        return nil, err
    }
    
    // Return response
    return &types.GenericAPIResponse{
        Success: true,
        Values: result,
    }, nil
}
```

### Java Controller
```java
@RestController
@RequestMapping("/api/v1")
public class AuthController {
    
    @Autowired
    private AuthService authService;
    
    @PostMapping("/login")
    public ResponseEntity<GenericResponse<LoginResponse>> login(@RequestBody LoginRequest request) {
        try {
            LoginResponse result = authService.login(request.getEmail(), request.getPassword());
            return ResponseEntity.ok(new GenericResponse<>(true, result));
        } catch (Exception e) {
            return ResponseEntity.badRequest()
                .body(new GenericResponse<>(false, null, e.getMessage()));
        }
    }
}
```

**Key Differences:**
- **Go**: Manual HTTP handling (`http.ResponseWriter`, `*http.Request`)
- **Java**: Spring handles HTTP automatically (`@RequestBody`, `ResponseEntity`)
- **Go**: Explicit error returns (`error` return type)
- **Java**: Exceptions thrown and handled by `@ExceptionHandler`
- **Go**: Manual dependency injection (constructor injection)
- **Java**: `@Autowired` for dependency injection

---

## 5. Services (Business Logic)

### Go Service (`core/services/auth_service.go`)
```go
type authService struct {
    store model.DataStore
}

func (s *authService) Login(ctx context.Context, email, password string) (*types.LoginResponse, error) {
    // Get user from database
    user, err := s.store.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, fmt.Errorf("invalid credentials")
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, fmt.Errorf("invalid credentials")
    }
    
    // Generate JWT
    token, err := s.generateJWT(user.ID, user.Email)
    if err != nil {
        return nil, err
    }
    
    return &types.LoginResponse{
        Token: token,
        User: types.UserInfo{...},
    }, nil
}
```

### Java Service
```java
@Service
public class AuthService {
    
    @Autowired
    private UserRepository userRepository;
    
    @Autowired
    private JwtTokenProvider jwtTokenProvider;
    
    public LoginResponse login(String email, String password) {
        // Get user from database
        User user = userRepository.findByEmail(email)
            .orElseThrow(() -> new InvalidCredentialsException("Invalid credentials"));
        
        // Verify password
        if (!passwordEncoder.matches(password, user.getPasswordHash())) {
            throw new InvalidCredentialsException("Invalid credentials");
        }
        
        // Generate JWT
        String token = jwtTokenProvider.generateToken(user.getId(), user.getEmail());
        
        return new LoginResponse(token, new UserInfo(user));
    }
}
```

**Key Differences:**
- **Go**: Explicit error returns (`error` type)
- **Java**: Exceptions (`throw new InvalidCredentialsException()`)
- **Go**: Context propagation (`context.Context`)
- **Java**: No explicit context (Spring handles it)
- **Go**: Interface-based design (explicit interfaces)
- **Java**: Interface + implementation pattern

---

## 6. Repository/DataStore Layer

### Go DataStore (`datastore/store.go`)
```go
type dataStore struct {
    db *gorm.DB
}

func (ds *dataStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := ds.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (ds *dataStore) GetSeatByIDForUpdate(ctx context.Context, seatID uint) (*model.ShowSeat, error) {
    var seat model.ShowSeat
    err := ds.db.WithContext(ctx).
        Set("gorm:query_option", "FOR UPDATE").
        Where("id = ?", seatID).
        First(&seat).Error
    return &seat, err
}
```

### Java Repository
```java
@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);
}

@Repository
public class SeatRepository {
    
    @PersistenceContext
    private EntityManager entityManager;
    
    public Optional<ShowSeat> findByIdForUpdate(Long seatId) {
        return entityManager.createQuery(
            "SELECT s FROM ShowSeat s WHERE s.id = :id",
            ShowSeat.class
        )
        .setParameter("id", seatId)
        .setLockMode(LockModeType.PESSIMISTIC_WRITE)  // FOR UPDATE
        .getResultStream()
        .findFirst();
    }
}
```

**Key Differences:**
- **Go**: GORM (Go ORM) with method chaining
- **Java**: JPA/Hibernate with annotations or JPQL
- **Go**: `FOR UPDATE` via query option
- **Java**: `LockModeType.PESSIMISTIC_WRITE` for row locking
- **Go**: Explicit context for cancellation/timeout
- **Java**: `@Transactional` for transaction management

---

## 7. Transactions

### Go Transaction (`datastore/store.go`)
```go
func (ds *dataStore) Begin(ctx context.Context) (model.Transaction, error) {
    tx := ds.db.WithContext(ctx).Begin()
    return &transaction{tx: tx}, nil
}

// In service layer
tx, err := store.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback()

seat, err := store.GetSeatByIDForUpdate(ctx, seatID)
if err != nil {
    return err
}

seat.Status = "LOCKED"
if err := store.UpdateSeat(ctx, seat); err != nil {
    return err
}

return tx.Commit()
```

### Java Transaction
```java
@Service
@Transactional
public class SeatService {
    
    @Autowired
    private SeatRepository seatRepository;
    
    public void lockSeat(Long seatId) {
        // Transaction automatically managed by @Transactional
        ShowSeat seat = seatRepository.findByIdForUpdate(seatId)
            .orElseThrow(() -> new SeatNotFoundException());
        
        seat.setStatus("LOCKED");
        seatRepository.save(seat);
        // Auto-commit on method exit, auto-rollback on exception
    }
}
```

**Key Differences:**
- **Go**: Manual transaction management (`Begin()`, `Commit()`, `Rollback()`)
- **Java**: Declarative transactions (`@Transactional`)
- **Go**: Explicit `defer` for cleanup
- **Java**: Automatic rollback on exceptions

---

## 8. Dependency Injection

### Go DI (`cmd/main.go`)
```go
// Manual dependency injection
store := datastore.NewDataStore(db)
authService := services.NewAuthService(clients, store)
ctrl := controllers.NewController(authService, ...)
```

### Java DI
```java
@SpringBootApplication
public class Application {
    // Spring automatically injects dependencies
}

@Service
public class AuthService {
    @Autowired
    private UserRepository userRepository;  // Auto-injected
}
```

**Key Differences:**
- **Go**: Constructor injection (manual wiring)
- **Java**: Field/constructor injection via `@Autowired`
- **Go**: Explicit initialization in `main()`
- **Java**: Spring container manages lifecycle

---

## 9. Error Handling

### Go Error Handling
```go
func (s *authService) Login(...) (*types.LoginResponse, error) {
    user, err := s.store.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, fmt.Errorf("invalid credentials: %w", err)  // Wrap error
    }
    return result, nil
}

// Caller checks error
result, err := authService.Login(ctx, email, password)
if err != nil {
    return nil, err  // Propagate up
}
```

### Java Error Handling
```java
public LoginResponse login(String email, String password) {
    User user = userRepository.findByEmail(email)
        .orElseThrow(() -> new InvalidCredentialsException("Invalid credentials"));
    return result;
}

// Caller catches exception
try {
    LoginResponse result = authService.login(email, password);
} catch (InvalidCredentialsException e) {
    // Handle error
}
```

**Key Differences:**
- **Go**: Explicit error returns (part of function signature)
- **Java**: Exceptions (checked/unchecked)
- **Go**: Error wrapping with context (`fmt.Errorf("...: %w", err)`)
- **Java**: Exception chaining (`throw new Exception("...", cause)`)

---

## 10. Interfaces

### Go Interface (`core/services/interfaces.go`)
```go
type AuthServiceInterface interface {
    Login(ctx context.Context, email, password string) (*types.LoginResponse, error)
}

// Implementation
type authService struct {
    store model.DataStore
}

func NewAuthService(clients *coretypes.Clients, store model.DataStore) AuthServiceInterface {
    return &authService{store: store}
}
```

### Java Interface
```java
public interface AuthService {
    LoginResponse login(String email, String password);
}

// Implementation
@Service
public class AuthServiceImpl implements AuthService {
    @Override
    public LoginResponse login(String email, String password) {
        // Implementation
    }
}
```

**Key Differences:**
- **Go**: Interfaces are implicit (if type has methods, it implements interface)
- **Java**: Explicit `implements` keyword
- **Go**: Interface naming convention: `XxxInterface` or just `Xxx`
- **Java**: Interface + `Impl` suffix for implementation

---

## 11. HTTP Routing

### Go Router (`api/v1/router.go`)
```go
routes := []Route{
    {
        Path:         "/api/v1/login",
        RequestMethod: http.MethodPost,
        Handler:      controllers.ResponseHandler(ctrl.LoginHandler),
        SkipAuth:     true,
    },
}

for _, route := range routes {
    router.Handle(route.Path, handler).Methods(route.RequestMethod)
}
```

### Java Router
```java
@RestController
@RequestMapping("/api/v1")
public class AuthController {
    
    @PostMapping("/login")
    public ResponseEntity<...> login(...) {
        // Handler
    }
}
```

**Key Differences:**
- **Go**: Manual route registration with Gorilla Mux
- **Java**: Annotation-based routing (`@PostMapping`, `@GetMapping`)
- **Go**: Explicit middleware/interceptor chain
- **Java**: `@PreAuthorize`, `@Valid` annotations

---

## 12. Middleware/Interceptors

### Go Interceptor (`api/v1/interceptors/interceptors.go`)
```go
func AuthInterceptor(errorHandler func(error, http.ResponseWriter, *http.Request)) Interceptor {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            token := extractToken(r)
            if token == "" {
                errorHandler(errors.New("unauthorized"), w, r)
                return
            }
            // Validate token, set user ID in context
            next(w, r)
        }
    }
}
```

### Java Interceptor
```java
@Component
public class AuthInterceptor implements HandlerInterceptor {
    
    @Override
    public boolean preHandle(HttpServletRequest request, 
                           HttpServletResponse response, 
                           Object handler) {
        String token = extractToken(request);
        if (token == null || !validateToken(token)) {
            response.setStatus(HttpStatus.UNAUTHORIZED.value());
            return false;
        }
        return true;
    }
}
```

**Key Differences:**
- **Go**: Function-based middleware (higher-order functions)
- **Java**: Interface-based (`HandlerInterceptor`)
- **Go**: Manual chain building
- **Java**: Spring configuration for interceptor chain

---

## Summary: Key Concepts Mapped

| Go Concept | Java Equivalent | Notes |
|------------|------------------|-------|
| `struct` | `class` | Go structs are like classes without methods (methods are separate) |
| `interface` | `interface` | Go interfaces are implicit (duck typing) |
| `error` return type | `Exception` | Go uses explicit error returns |
| `context.Context` | `@Transactional` | Go uses context for cancellation/timeout |
| GORM | JPA/Hibernate | Both are ORMs |
| `defer` | `finally` | Go's defer is like finally but cleaner |
| Manual DI | `@Autowired` | Go requires manual wiring |
| `http.HandlerFunc` | `@RequestMapping` | Go uses functions, Java uses annotations |
| `SELECT ... FOR UPDATE` | `LockModeType.PESSIMISTIC_WRITE` | Both for row-level locking |

---

## Architecture Pattern Comparison

### Go: Three-Layer Architecture (Explicit)
```
HTTP Request
    ↓
Router (api/v1/router.go)
    ↓
Controller (api/v1/controllers/)
    ↓
Service (core/services/)
    ↓
DataStore (datastore/)
    ↓
Database
```

### Java: Spring MVC (Framework-Managed)
```
HTTP Request
    ↓
DispatcherServlet (Spring)
    ↓
@Controller (annotated)
    ↓
@Service (annotated)
    ↓
@Repository (annotated)
    ↓
Database
```

**Main Difference:**
- **Go**: You explicitly wire everything together
- **Java**: Spring Framework handles most wiring automatically

---

This should help you understand the Go backend by relating it to Java concepts! Ask me about any specific part you want to dive deeper into.
