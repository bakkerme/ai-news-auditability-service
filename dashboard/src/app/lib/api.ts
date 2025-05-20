export interface RunMetadata {
  id: string;
  personaName: string;
  runDate: string; // ISO 8601 date-time string
  overallModelUsed?: string;
  totalItems?: number;
  hasBenchmark?: boolean;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/v1'; // Default to local dev server

export async function getLatestRun(): Promise<RunMetadata> {
  try {
    const response = await fetch(`${API_BASE_URL}/runs/latest`);
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Unknown error fetching latest run' }));
      throw new Error(`Failed to fetch latest run: ${response.status} ${response.statusText} - ${errorData.message || JSON.stringify(errorData)}`);
    }
    const data: RunMetadata = await response.json();
    return data;
  } catch (error) {
    console.error('Error in getLatestRun:', error);
    // In a real app, you might want to throw a custom error or handle it differently
    // For now, rethrow to let the component handle it.
    throw error;
  }
} 