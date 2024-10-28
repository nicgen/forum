# MCD (Méthode Conceptuelle de Données)

### Entité : User
- **Attributs:** Id (PK), UUID, Username (Unique), Password, CreatedAt, Role, Email (Unique)
- **Relations:**  
1,n avec Posts (Un utilisateur peut créer plusieurs posts)  
1,n avec Comments (Un utilisateur peut commenter plusieurs posts)  
1,n avec ReportPost (Un utilisateur peut signaler plusieurs posts)  
1,n avec LikePost (Un utilisateur peut aimer/désaimer plusieurs posts)  
1,n avec ReportComment (Un utilisateur peut signaler plusieurs commentaires)  
1,n avec LikeComment (Un utilisateur peut aimer/désaimer plusieurs commentaires)  
1,1 avec RequestMod (Un utilisateur peut faire une seule demande de modération)  

### Entité : Posts

- **Attributs:** Id (PK), User_id (FK), Texte, CreatedAt, Title, UpdatedAt  
- **Relations:**  
n,1 avec User (Plusieurs posts peuvent être créés par un utilisateur)  
n,n avec Categories via Post_Categories (Un post peut être lié à plusieurs catégories)  
1,n avec Comments (Un post peut avoir plusieurs commentaires)  
1,n avec ReportPost (Un post peut être signalé plusieurs fois)  
1,n avec LikePost (Un post peut être aimé/désaimé plusieurs fois)  

### Entité : Categories

- **Attributs:** Id (PK), Name  
- **Relations:**  
n,n avec Posts via Post_Categories (Une catégorie peut regrouper plusieurs posts)  

### Entité : Post_Categories (Association)

- **Attributs:** Post_id (PK, FK), Categorie_id (PK, FK)  
- **Relations:**  
n,1 avec Posts  
n,1 avec Categories  

### Entité : Comments

- **Attributs:** Id (PK), Post_id (FK), User_id (FK), Texte, CreatedAt  
- **Relations:**  
n,1 avec Posts (Plusieurs commentaires peuvent être liés à un post)  
n,1 avec User (Plusieurs commentaires peuvent être créés par un utilisateur)  
1,n avec ReportComment (Un commentaire peut être signalé plusieurs fois)  
1,n avec LikeComment (Un commentaire peut être aimé/désaimé plusieurs fois)  

### Entité : ReportPost (Association)

- **Attributs:** Post_id (PK, FK), User_id (PK, FK)  
- **Relations:**  
n,1 avec Posts  
n,1 avec User  

### Entité : LikePost (Association)

- **Attributs:** Post_id (PK, FK), User_id (PK, FK), IsLike  
- **Relations:**  
n,1 avec Posts  
n,1 avec User  

### Entité : ReportComment (Association)

- **Attributs:** Post_id (PK, FK), Comment_id (PK, FK), User_id (PK, FK)  
- **Relations:**
n,1 avec Posts  
n,1 avec Comments  
n,1 avec User  

### Entité : LikeComment (Association)

- **Attributs:** Comment_id (PK, FK), User_id (PK, FK), IsLike  
- **Relations:**  
n,1 avec Comments  
n,1 avec User  

### Entité : RequestMod

- **Attributs:** User_id (PK, FK), Reason  
- **Relations:**  
1,1 avec User (Un utilisateur peut faire une seule demande de modération)  
