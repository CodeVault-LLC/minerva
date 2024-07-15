import { Category, SubCategory } from "~/lib/constants";
import type { RegexPattern } from "~/lib/types";

export const SECRET_PATTERNS: RegexPattern[] = [
  {
    pattern: /(?<=secret\s*:\s*['"])[^'"]+/gi,
    name: "Sensitive Data Keywords",
    description:
      "Matches common keywords related to sensitive data like keys, tokens, secrets, etc.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.TOKEN,
  },
  {
    // Google API key
    pattern: /AIza[0-9A-Za-z_-]{35}/gi,
    name: "Google API Key",
    description: "Matches Google API keys.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.API_KEY,
  },
  {
    // AWS Access Key ID
    pattern: /AKIA[0-9A-Z]{16}/gi,
    name: "AWS Access Key ID",
    description: "Matches AWS access key IDs.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.API_KEY,
  },
  {
    // Google Analytics Tracking ID
    pattern: /UA-[0-9]+-[0-9]+/gi,
    name: "Google Analytics Tracking ID",
    description: "Matches Google Analytics tracking IDs.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.TRACKING_ID,
  },

  // Firebase config object

  {
    pattern: /https?:\/\/[^\s/$.?#].[^\s"]*/gi,
    name: "URL Matching",
    description: "Matches HTTP and HTTPS URLs.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },
  {
    pattern: /http?:\/\/[^\s/$.?#].[^\s"]*/gi,
    name: "HTTP URL Matching",
    description: "Matches HTTP URLs.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },
  {
    pattern: /https?:\/\/[^\s/$.?#].[^\s"]*/gi,
    name: "HTTPS URL Matching",
    description: "Matches HTTPS URLs.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },
  {
    pattern: /ftp?:\/\/[^\s/$.?#].[^\s"]*/gi,
    name: "FTP URL Matching",
    description: "Matches FTP URLs.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },
  {
    pattern: /ssh?:\/\/[^\s/$.?#].[^\s"]*/gi,
    name: "SSH URL Matching",
    description: "Matches SSH URLs.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },
  {
    pattern: /sftp?:\/\/[^\s/$.?#].[^\s"]*/gi,
    name: "SFTP URL Matching",
    description: "Matches SFTP URLs.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },
  {
    pattern: /file:\/\/[^\s/$.?#].[^\s"]*/gi,
    name: "File URL Matching",
    description: "Matches file URLs.",
    category: Category.URL_PATH,
    subCategory: SubCategory.URL,
  },
  {
    pattern: /data:image\/[a-zA-Z]*;base64,[^'"\s]*/gi,
    name: "Base64 Image Data",
    description: "Matches base64 image data.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.IMAGE,
  },
  {
    pattern: /data:application\/[a-zA-Z]*;base64,[^'"\s]*/gi,
    name: "Base64 Application Data",
    description: "Matches base64 application data.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.APPLICATION,
  },
  {
    pattern: /data:audio\/[a-zA-Z]*;base64,[^'"\s]*/gi,
    name: "Base64 Audio Data",
    description: "Matches base64 audio data.",
    category: Category.SENSITIVE_DATA,
    subCategory: SubCategory.AUDIO,
  },
];
