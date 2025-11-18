
#  Projet Gomoku

##  Objectif du Projet

L'objectif principal est de créer une **IA capable de battre des joueurs humains au Gomoku**.

##  Exigences Générales et Techniques

###  Spécifications du Programme

* **Executable Name:** Le programme exécutable doit être nommé `Gomoku`.
* **Langage/Librairies:** Vous êtes libre d'utiliser le langage et la librairie d'interface graphique de votre choix.
* **Stabilité:** Le programme ne doit **jamais planter** (même par manque de mémoire) ou quitter inopinément. Un crash entraînera une note de 0.
* **Makefile:** Un `Makefile` est obligatoire et doit produire votre programme sans relink. Il doit contenir au minimum les règles : `$(NAME)`, `all`, `clean`, `fclean` et `re`.
* **Performance de l'IA:**
    * Le temps moyen pour trouver un coup ne doit **pas dépasser une demi-seconde**.
    * L'implémentation doit rester performante ; une IA qui gagne trop lentement ou qui semble "paresseuse" (faible profondeur de recherche, implémentation naïve) ne recevra pas tous les points.
* **Interface Graphique (GUI):** Une interface graphique utilisable et "vaguement agréable" à l'œil est obligatoire.

###  Exigence de l'Interface

L'interface utilisateur doit afficher un **chronomètre** qui indique le temps que met l'IA pour trouver son coup suivant. *L'absence de ce chronomètre invalide le projet*.

##  Règles du Jeu (Additionnelles au Gomoku standard)

Le Gomoku se joue sur un **Goban de 19x19**.

* **Condition de Victoire par Alignement:** Un joueur gagne en alignant **cinq pierres ou plus** de sa couleur.
* **Capture (Ninuki-renju ou Pente):**
    * Vous pouvez retirer une **paire** de pierres de l'adversaire en les flanquant avec vos propres pierres.
    * Les captures ne peuvent se faire que par paires (pas de pierre unique, pas plus de 2 pierres d'affilée).
    * **Victoire par Capture:** Réaliser **dix captures** (cinq paires) met fin à la partie et donne la victoire.
* **Capture en Fin de Partie (Endgame Capture):**
    * Une victoire par alignement de cinq pierres est validée **uniquement si l'adversaire ne peut pas briser cette ligne en capturant une paire**.
    * Si le joueur gagnant a déjà perdu quatre paires, et que l'adversaire peut capturer une cinquième, l'adversaire gagne par capture.
* **Interdiction des Doubles-Trois (Free-threes):**
    * Un **Triple-Libre (Free-three)** est un alignement de trois pierres qui, si non bloqué, permet un alignement de quatre pierres indéfendable (avec deux extrémités libres).
    * Il est **interdit** de jouer un coup qui crée **deux Triples-Libres simultanés** (un Double-Trois).
    * *Note: Il est permis d'introduire un Double-Trois en réalisant une capture de paire*.

##  Partie Obligatoire (AI)

Vous devez implémenter la possibilité de jouer dans les modes suivants:

* **Contre votre programme (IA vs Humain):** L'IA doit être capable de gagner la partie sans que l'humain la laisse gagner, et elle doit adapter sa stratégie aux coups de l'adversaire.
* **Contre un autre joueur humain (Hotseat):** Avec une **fonction de suggestion de coup** fournie par l'IA.

###  Algorithme de l'IA

* **Algorithme:** Vous devez utiliser l'algorithme **Min-Max** (ou une variante).
* **Profondeur de Recherche:** Pour valider le projet, votre IA doit rechercher **au minimum 10 niveaux** de profondeur dans son arbre de jeu. Une profondeur inférieure ne permettra pas d'atteindre la note maximale.
* **Fonction Heuristique:**
    * L'IA doit générer un arbre de solutions et choisir le meilleur coup selon cet arbre.
    * Vous devez concevoir une fonction heuristique **efficace et rapide** pour évaluer la valeur des nœuds terminaux. L'heuristique est considérée comme la partie la plus difficile du projet.
* **Outil de Débogage:** Il est fortement recommandé d'implémenter un processus de débogage permettant d'examiner le raisonnement de l'IA pendant son exécution.

##  Partie Bonus

La partie bonus ne sera évaluée **QUE SI la partie obligatoire est PARFAITE** (intégralement faite et fonctionnelle sans aucun dysfonctionnement).

* **Bonus Suggéré:** Implémenter la possibilité de choisir les règles du jeu au début d'une partie (par exemple, les conditions de départ comme Standard, Pro, Swap, Swap2, etc.).
* **Autres Bonus:** Tout autre bonus intéressant ou utile sera pris en compte (jusqu'à 5 bonus maximum).

##  Défense du Projet

Lors de la session de défense, vous devez être prêt à:

* **Expliquer en détail votre implémentation de l'algorithme Minimax** (ou de la variante choisie).
* **Expliquer en détail votre fonction heuristique** (elle doit être précise, rapide, et vous devez la comprendre parfaitement).
* Démontrer que vous avez correctement implémenté **toutes les règles du jeu** imposées par le sujet.
* Faire tourner votre programme contre l'évaluateur.
