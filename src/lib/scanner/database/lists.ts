import { z } from "zod";

export const ClientDatabases = z.enum([
  "firebase",
  "aws",
  "google-cloud",
  "azure",
]);

export const regexes = [
  {
    name: ClientDatabases.Enum.firebase,
    patterns: [
      // Import statements
      /\bimport\s+.*\bfirebase\b.*from\s+['"]firebase['"]/g,
      /\brequire\s*\(\s*['"]firebase['"]\s*\)/g,
      // Initialization or usage patterns
      /\bfirebase\.initializeApp\b\s*\(/g,
      /\bfirebase\.auth\b\s*\(/g,
      /\bfirebase\.database\b\s*\(/g,
      /\bfirebase\.firestore\b\s*\(/g,
      // URL references
      /(?:^|[\W])firebaseio\.com(?:$|[\W])/g,
      /(?:^|[\W])firebaseio\.com\/[\w.-]+(?:$|[\W])/g,
      /(?:^|[\W])firebasedatabase\.app(?:$|[\W])/g,
      /(?:^|[\W])firebasedatabase\.app\/[\w.-]+(?:$|[\W])/g,
      /(?:^|[\W])firebaseapp\.com(?:$|[\W])/g,
      /(?:^|[\W])firebaseapp\.com\/[\w.-]+(?:$|[\W])/g,
    ],
  },
  {
    name: ClientDatabases.Enum.aws,
    patterns: [
      // Import statements
      /\bimport\s+.*\baws-sdk\b.*from\s+['"]aws-sdk['"]/g,
      /\brequire\s*\(\s*['"]aws-sdk['"]\s*\)/g,
      // Initialization or usage patterns
      /\bAWS\.config\.update\b\s*\(/g,
      /\bnew\s+AWS\.[A-Z][a-zA-Z]+\b\s*\(/g,
      // URL references
      /(?:^|[\W])dynamodb\.[\w.-]+\.amazonaws\.com(?:$|[\W])/g,
      /(?:^|[\W])rds\.[\w.-]+\.amazonaws\.com(?:$|[\W])/g,
      /(?:^|[\W])s3\.[\w.-]+\.amazonaws\.com(?:$|[\W])/g,
      /(?:^|[\W])awsapps\.com(?:$|[\W])/g,
      /(?:^|[\W])cloudfront\.net(?:$|[\W])/g,
    ],
  },
  {
    name: ClientDatabases.Enum["google-cloud"],
    patterns: [
      // Import statements
      /\bimport\s+.*\b@google-cloud\/firestore\b.*from\s+['"]@google-cloud\/firestore['"]/g,
      /\brequire\s*\(\s*['"]@google-cloud\/firestore['"]\s*\)/g,
      // Initialization or usage patterns
      /\bnew\s+Firestore\b\s*\(/g,
      /\bnew\s+Datastore\b\s*\(/g,
      /\bnew\s+Storage\b\s*\(/g,
      // URL references
      /(?:^|[\W])datastore\.googleapis\.com(?:$|[\W])/g,
      /(?:^|[\W])firestore\.googleapis\.com(?:$|[\W])/g,
      /(?:^|[\W])storage\.googleapis\.com(?:$|[\W])/g,
      /(?:^|[\W])appspot\.com(?:$|[\W])/g,
    ],
  },
  {
    name: ClientDatabases.Enum.azure,
    patterns: [
      // Import statements
      /\bimport\s+.*\b@azure\/cosmos\b.*from\s+['"]@azure\/cosmos['"]/g,
      /\brequire\s*\(\s*['"]@azure\/cosmos['"]\s*\)/g,
      // Initialization or usage patterns
      /\bnew\s+CosmosClient\b\s*\(/g,
      /\bnew\s+BlobServiceClient\b\s*\(/g,
      /\bnew\s+TableServiceClient\b\s*\(/g,
      // URL references
      /(?:^|[\W])database\.windows\.net(?:$|[\W])/g,
      /(?:^|[\W])documents\.azure\.com(?:$|[\W])/g,
      /(?:^|[\W])table\.core\.windows\.net(?:$|[\W])/g,
      /(?:^|[\W])blob\.core\.windows\.net(?:$|[\W])/g,
      /(?:^|[\W])azurewebsites\.net(?:$|[\W])/g,
    ],
  },
];

export const loginRegexes = [
  // Firebase Auth Config (ApiKey, AuthDomain, ProjectId etc.)
  {
    name: "Firebase Auth Config",
    patterns: [
      /apiKey\s*:\s*['"`][A-Za-z0-9_-]+['"`]/g,
      /authDomain\s*:\s*['"`][A-Za-z0-9_\-.]+['"`]/g,
      /databaseURL\s*:\s*['"`][A-Za-z0-9_\-:./]+['"`]/g,
      /projectId\s*:\s*['"`][A-Za-z0-9_-]+['"`]/g,
      /storageBucket\s*:\s*['"`][A-Za-z0-9_\-.]+['"`]/g,
      /messagingSenderId\s*:\s*['"`][0-9]+['"`]/g,
      /appId\s*:\s*['"`][A-Za-z0-9_\-:.]+['"`]/g,
      /measurementId\s*:\s*['"`][A-Za-z0-9_-]+['"`]/g,
    ],
  },
  // AWS Cognito Config (UserPoolId, ClientId etc.)
  {
    name: "AWS Cognito Config",
    patterns: [
      /userPoolId:\s*['"][\w-]+['"]/g,
      /userPoolWebClientId:\s*['"][\w-]+['"]/g,
    ],
  },

  // Google Auth Config (ClientId, ApiKey etc.)
  {
    name: "Google Auth Config",
    patterns: [/clientId:\s*['"][\w-]+['"]/g, /apiKey:\s*['"][\w-]+['"]/g],
  },

  // Azure AD B2C Config (ClientId, Authority etc.)
  {
    name: "Azure AD B2C Config",
    patterns: [/clientId:\s*['"][\w-]+['"]/g, /authority:\s*['"][\w-]+['"]/g],
  },

  // GitHub Auth Config (ClientId, ClientSecret etc.)
  {
    name: "GitHub Auth Config",
    patterns: [
      /clientId:\s*['"][\w-]+['"]/g,
      /clientSecret:\s*['"][\w-]+['"]/g,
    ],
  },
];
