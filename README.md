# Auto-Syntax-Fixer

Un outil CLI puissant et minimaliste qui corrige automatiquement les erreurs de syntaxe dans les repositories GitHub.

## 🚀 Fonctionnalités

- **Clone automatique** : Récupère n'importe quel repo GitHub (public/privé)
- **Détection multi-langages** : Go, Python, JavaScript, TypeScript
- **Corrections automatiques** : Syntaxe, formatage, imports manquants
- **Workflow Git intégré** : Crée une branche, commit et push automatiquement
- **Rapport détaillé** : Résumé des corrections appliquées

## 📦 Installation

```bash
git clone https://github.com/votre-username/auto-syntax-fixer
cd auto-syntax-fixer
go build -o auto-syntax-fixer
```

## 🔧 Usage

### Commande de base
```bash
./auto-syntax-fixer --repo https://github.com/username/repository
```

### Avec authentification (repos privés)
```bash
./auto-syntax-fixer --repo https://github.com/username/private-repo --token YOUR_GITHUB_TOKEN
```

### Options disponibles
```bash
--repo       URL du repository GitHub (requis)
--token      Token GitHub pour les repos privés
--branch     Branche à utiliser (défaut: main)
--dry-run    Mode test sans modifications
--output     Fichier de sortie pour le rapport
```

## 🛠️ Types de corrections

### Go
- Formatage avec `go fmt`
- Imports manquants avec `goimports`
- Corrections syntaxiques basiques

### Python
- Formatage avec `autopep8`
- Tri des imports avec `isort`
- Corrections d'indentation et syntaxe

### JavaScript/TypeScript
- Formatage avec `prettier`
- Corrections des points-virgules manquants
- Fixes d'imports ES6

## 📋 Exemples

### Correction d'un repo Go
```bash
./auto-syntax-fixer --repo https://github.com/example/go-project
```

### Mode dry-run (test sans modifications)
```bash
./auto-syntax-fixer --repo https://github.com/example/project --dry-run
```

### Générer un rapport
```bash
./auto-syntax-fixer --repo https://github.com/example/project --output report.txt
```

## 🏗️ Architecture

```
auto-syntax-fixer/
├── main.go              # CLI et orchestration
├── fixer/
│   ├── detector.go      # Détection des langages
│   ├── go.go           # Corrections Go
│   ├── python.go       # Corrections Python
│   └── javascript.go   # Corrections JS/TS
├── git/
│   └── operations.go    # Opérations Git
└── go.mod
```

## 🔒 Sécurité

- Validation des inputs
- Gestion sécurisée des tokens GitHub
- Isolation des opérations dans des dossiers temporaires

## 🤝 Contribution

1. Fork le projet
2. Créez une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit vos changements (`git commit -m 'Add AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrez une Pull Request

## 📝 License

Distribué sous licence MIT. Voir `LICENSE` pour plus d'informations.

## 🚨 Limitations

- Focus sur les corrections syntaxiques uniquement
- Pas de refactoring complexe
- Pas de correction de logique métier
- Nécessite Git installé sur le système

## 🔧 Dépendances

- Go 1.21+
- Git
- Outils optionnels : `goimports`, `autopep8`, `isort`, `prettier`
