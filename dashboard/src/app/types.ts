export interface Item {
  id: string;
  title: string;
  summary: string;
  commentSummary?: string;
  imageDescription?: string;
  webContentSummary?: string;
  link?: string;
  thumbnailUrl?: string;
  isRelevant: boolean;
  entry: Entry;
}

export interface Entry {
  title: string;
  link: Link;
  id: string;
  published: string;
  content: string;
  comments: EntryComment[];
  externalURLs: string[];
  imageURLs: string[];
  mediaThumbnail: MediaThumbnail;
  imageDescription: string;
  webContentSummaries: string[];
}

export interface EntryComment {
  content: string;
}

export interface Link {
  href: string;
}

export interface MediaThumbnail {
  url: string;
}

export interface EntrySummary {
  rawInput: string;
  results: Item;
  processingTimeMs: number;
}

export interface Persona {
  name: string;
  feedUrl: string; // URL of the RSS feed (e.g., "https://reddit.com/r/localllama.rss")
  topic: string; // Main subject area (e.g., "AI Technology", "Gardening")
  personaIdentity: string; // Core identity and expertise of the persona
  basePromptTask: string; // Task description for individual item analysis
  summaryPromptTask: string; // Task description for summary generation
  focusAreas: string[]; // List of topics/keywords to prioritize
  relevanceCriteria: string[]; // List of criteria for relevance analysis
  summaryAnalysis: string[]; // Focus areas for summary analysis
  exclusionCriteria: string[]; // List of criteria to explicitly exclude items
}

export interface ImageSummary {
  imageURL: string;
  imageDescription: string;
  title?: string;
  entryID?: string;
  processingTimeMs: number;
}

export interface WebContentSummary {
  url: string;
  originalContent: string;
  summary: string;
  title?: string;
  entryID?: string;
  processingTimeMs: number;
}

// KeyDevelopment represents a key development and its referenced item
export interface KeyDevelopment {
  text: string;
  itemID: string;
}

// SummaryResponse represents an overall summary of multiple relevant AI news items
export interface SummaryResponse {
  keyDevelopments: KeyDevelopment[];
}

export interface RunData {
  runId: string;
  entrySummaries: EntrySummary[];
  imageSummaries: ImageSummary[];
  webContentSummaries: WebContentSummary[];
  overallSummary: SummaryResponse; // Added to sync with Go model
  persona: Persona;
  runDate: string; // Assuming ISO string format from Go's time.Time
  overallModelUsed: string;
  imageModelUsed: string;
  webContentModelUsed: string;
  totalProcessingTime: number;
  entryTotalProcessingTime: number;
  imageTotalProcessingTime: number;
  webContentTotalProcessingTime: number;
  successRate?: number;
} 