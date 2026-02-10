import { Routes } from '@angular/router';
import { HomeComponent } from './features/home/home';
import { GameComponent } from './features/game/game';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'game', component: GameComponent }
];