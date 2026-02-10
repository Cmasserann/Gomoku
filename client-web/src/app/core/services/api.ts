import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';

export interface ServerStatus {
  goban_free: boolean;
  pending_invitation: boolean;
}

export interface GameInitResponse {
  session_token: string;
}

export interface CreateGameResponse {
  session_token: string;
  invitation_token: string;
}

export interface MoveResponse {
  x: number;
  y: number;
  status: string; // ex: "ongoing", "win", "draw"
}

@Injectable({ providedIn: 'root' })
export class ApiService {
  private http = inject(HttpClient);
  private readonly API_URL = 'http://localhost:8080';

  sessionToken = signal<string | null>(null);

  startAIGame() {
    return this.http.post<GameInitResponse>(`${this.API_URL}/create`, { "ai_mode": true, "local_mode": false });
  }

  getAIMove() {
    const token = this.sessionToken();
    // On passe le token dans les headers ou en paramètre selon ton API Go
    return this.http.get<MoveResponse>(`${this.API_URL}/ai-suggest`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
  }

  createGame() {
  // On part du principe que ton backend a un endpoint POST /create
  return this.http.post<CreateGameResponse>(`${this.API_URL}/create`, {});
  }

  // STOCKAGE : On initialise à 'null' (en attente de réponse)
  status = signal<ServerStatus | null>(null);

  checkStatus() {
    // REQUÊTE : On appelle l'URL /status
    this.http.get<ServerStatus>(`${this.API_URL}/status`)
      .subscribe({
        // MISE À JOUR : Quand la donnée arrive, on l'injecte dans le signal
        next: (data) => this.status.set(data),
        // ERREUR : Si le serveur est éteint, on peut mettre un état spécifique
        error: () => this.status.set(null) 
      });
  }
}
