import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';

type Player = 'black' | 'white';
type Cell = Player | null;

@Component({
  selector: 'app-game',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './game.html',
  styleUrl: './game.scss'
})
export class GameComponent {
  // On crée un plateau de 19x19 rempli de 'null'
  boardSize = 19;
  grid = signal<Cell[][]>(
    Array(this.boardSize).fill(null).map(() => Array(this.boardSize).fill(null))
  );

  currentPlayer = signal<Player>('black'); // Le noir commence toujours

  playMove(x: number, y: number) {
    const currentBoard = this.grid();
    
    // Si la case est déjà prise, on ne fait rien
    if (currentBoard[y][x] !== null) return;

    // On place le pion
    currentBoard[y][x] = this.currentPlayer();
    this.grid.set([...currentBoard]); // On met à jour le signal (shallow copy)

    // On change de joueur
    this.currentPlayer.set(this.currentPlayer() === 'black' ? 'white' : 'black');

    // TODO: Ici on appellera le backend Go pour l'IA ou la vérification de victoire
  }
}