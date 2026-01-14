
# 💡 API Documentation: Go Game Server

Cette API gère une instance unique de serveur de jeu de Go. Elle supporte les modes de jeu local, contre l'IA, ou en ligne via un système d'invitation.

## 🔑 Système d'Authentification

L'API utilise deux types de jetons (tokens) :

1. **Invitation Token (4 chars) :** Généré pour le joueur 2 lors de la création d'une partie multijoueur.
2. **Session Token (16 chars) :** Utilisé pour authentifier chaque coup (`/move`, `/ai-suggest`).

---

## 🛰️ État du Serveur

### 1. Vérifier la disponibilité

`GET /status`

Permet de savoir si une partie est déjà lancée avant d'essayer d'en créer une.

**Réponse (200 OK) :**

```json
{
  "goban_free": true
}

```

### 2. Récupérer le plateau

`GET /board`

Renvoie l'état complet de la partie en cours.

**Réponse (200 OK) :**

* `board`: Matrice 2D d'entiers (0: vide, 1: noir, 2: blanc).
* `captured_b` / `captured_w`: Nombre de pierres capturées.
* `turn`: Numéro du tour actuel.
* `goban_free`: État d'occupation du serveur.

---

## 🎮 Gestion des Sessions

### 3. Créer une partie

`POST /create`

**Request Body :**

```json
{
  "ai_mode": false,
  "local_mode": false
}

```

**Réponse (200 OK) :**

* `player_one`: Ton token de session (16 chars).
* `player_two`: Token d'invitation (4 chars) à donner à un ami, **OU** vide si IA/Local.

### 4. Abandon de la partie

`POST /giveUp`

Donne la victoire a l'adversere et libere la partie.

**Request Body :** `{"token": "Ton token de session (16 chars)"}`

**Réponse (200 OK) :**

* `"message": "Game over."`

### 5. Rejoindre une partie

`POST /join`

Échange un code d'invitation contre un token de session.

**Request Body :** `{"token": "a1b2"}`

**Réponse (200 OK) :**

* `token`: Ton token de session définitif (16 chars).

---

## 🕹️ Gameplay

### 6. Jouer un coup

`POST /move`

Soumet un coup au serveur. Si le mode IA est activé, l'IA répondra immédiatement dans la même requête.

**Request Body :**

```json
{
  "x": 10,
  "y": 5,
  "token": "ton_token_16_chars"
}

```

**Réponses :**

* **200 OK :** Renvoie le `board` et le `turn`.
* Si l'IA a joué : inclut un champ `time_us` (microsecondes).
* Si la partie est finie : inclut un champ `winner`.


* **401 Unauthorized :** Token invalide ou mauvais tour.
* **400 Bad Request :** Coup illégal (règles du Go).

### 7. Suggestion de l'IA

`POST /ai-suggest`

Demande à l'IA quel serait le meilleur coup sans le jouer.

**Request Body :** `{"token": "..."}`

**Réponse (200 OK) :**

```json
{
  "x": 5,
  "y": 5,
  "time_us": 1250
}

```

---

## 🛠️ Debug & Administration

### 8. Debug Mode

`POST /debug`

Permet de modifier l'état interne du serveur pour tester des situations spécifiques.

**Request Body :**

* `board`: (Optionnel) Injecte une matrice 2D.
* `captured_b` / `captured_w`: (Optionnel) Modifie les scores.
* `reset_board`: (bool) Si `true`, réinitialise tout et libère le serveur.
