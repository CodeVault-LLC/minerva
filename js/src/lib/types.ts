import type { Category, SubCategory } from "./constants";

export interface RegexPattern {
  pattern: RegExp;
  name: string;
  description: string;
  category: Category;
  subCategory?: SubCategory;
}

export interface Secret {
  name: string;
  results:
    | {
        match: string;
        line: number;
      }[]
    | never[];
}
