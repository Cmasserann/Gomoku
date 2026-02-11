import { Component, OnInit, inject, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from '../../core/services/api';
import { Router } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule], 
  templateUrl: './home.html',
  styleUrl: './home.scss'
})
export class HomeComponent implements OnInit {
  // 1. Injection du service
  private apiService = inject(ApiService);

  private router = inject(Router);

  // 2. Données pour le titre
  title = "五目並べ";
  letters = this.title.split('');

  // 3. Accès au signal de stockage du service
  status = this.apiService.status;

  // 4. Logique calculée (Signaux dérivés)
  // Ces variables se mettent à jour toutes seules dès que 'status' change
  canCreate = computed(() => this.status()?.goban_free === true);
  
  canJoin = computed(() => this.status()?.pending_invitation === true);
  
  isFull = computed(() => 
    this.status() !== null && 
    this.status()?.goban_free === false && 
    this.status()?.pending_invitation === false
  );

  ngOnInit() {
    // 5. Lancement de l'appel API au démarrage
    this.apiService.checkStatus();
  }

  // Fonctions pour les boutons (on les remplira après)
  startGame(mode: string) {
    if (mode === 'ai') {
      this.apiService.startAIGame().subscribe({
        next: (res) => {
          this.apiService.sessionToken.set(res.session_token);
          this.router.navigate(['/game']);
        },
        error: (err) => alert("Le serveur Go n'a pas pu créer la partie.")
      });
    }
  }

  createRoom() {
    console.log('Création de room...');
  }

  joinRoom() {
    console.log('Rejoint la room...');
  }
}