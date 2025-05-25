import Layout from './components/Layout';
import MetadataCard from './components/MetadataCard';
import AccordionItem from './components/AccordionItem';
import TruncatableContent from './components/TruncatableContent';
import { RunData } from './types'; // Import the RunData interface

async function getLatestRunData(): Promise<RunData> {
  try {
    // TODO: Make the API URL configurable, perhaps via environment variables
    const res = await fetch('http://localhost:8080/v1/runs/latest', {
      cache: 'no-store', // Ensure fresh data on each request
    });

    if (!res.ok) {
      const errorBody = await res.text();
      console.error(`Error fetching latest run data: ${res.status} ${res.statusText} - ${errorBody}`);
      // Throw an error to be caught by Next.js error handling
      throw new Error(`Failed to fetch latest run data: ${res.status} ${res.statusText}. Body: ${errorBody}`);
    }
    const data: RunData = await res.json();
    console.log('Latest run data:', data);
    return data;
  } catch (error) {
    // Log the error and re-throw to ensure Next.js handles it
    console.error('Exception during fetch latest run data:', error);
    if (error instanceof Error && error.message.startsWith('Failed to fetch')) {
      throw error; // Re-throw specific fetch errors
    }
    throw new Error(`An unexpected error occurred while fetching latest run data: ${error instanceof Error ? error.message : String(error)}`);
  }
}

export default async function HomePage() {
  let latestRun: RunData | null = null;
  let fetchError: string | null = null;

  try {
    latestRun = await getLatestRunData();
  } catch (error) {
    console.error("Error in HomePage while fetching data:", error);
    fetchError = error instanceof Error ? error.message : "An unknown error occurred.";
  }

  if (fetchError || !latestRun) {
    return (
      <Layout>
        <div className="container mx-auto p-4">
          <h1 className="text-2xl font-bold text-red-600 dark:text-red-400 mb-4">Error Loading Data</h1>
          <p className="text-red-500 dark:text-red-300 text-center text-xl">
            Failed to load latest run data. Please check if the service is running or try again later.
          </p>
          {fetchError && (
            <pre className="mt-4 p-4 bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-200 border border-red-300 dark:border-red-700 rounded-md whitespace-pre-wrap">
              {fetchError}
            </pre>
          )}
        </div>
      </Layout>
    );
  }
  
  const runId = latestRun.runId; // Use the actual run ID from the API response

  return (
    <Layout>
      <MetadataCard
        runId={runId}
        timestamp={latestRun.runDate}
        overallModel={latestRun.overallModelUsed}
        imageModel={latestRun.imageModelUsed}
        webContentModel={latestRun.webContentModelUsed}
        totalProcessingTime={latestRun.totalProcessingTime}
      />

      <div className="mb-6 p-6 bg-white dark:bg-slate-800 shadow-lg rounded-lg border border-slate-200 dark:border-slate-700">
        <h2 className="text-2xl font-semibold mb-4 text-slate-800 dark:text-slate-100 border-b dark:border-slate-600 pb-2">Persona Details</h2>
        <div className="space-y-4">
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Name</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">{latestRun.persona.name}</p>
          </div>
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Identity / System Prompt</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md whitespace-pre-wrap">{latestRun.persona.basePromptTask || "Not available"}</p>
          </div>
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Focus Areas</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">
              {Array.isArray(latestRun.persona.focusAreas) ? latestRun.persona.focusAreas.join(', ') : "Not available"}
            </p>
          </div>
           <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Exclusion Criteria</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">
              {Array.isArray(latestRun.persona.exclusionCriteria) ? latestRun.persona.exclusionCriteria.join(', ') : "Not available"}
            </p>
          </div>
          
        </div>
      </div>

      {/* Overall Summary Section */}
      {latestRun.overallSummary && latestRun.overallSummary.keyDevelopments && latestRun.overallSummary.keyDevelopments.length > 0 && (
        <div className="mb-6 p-6 bg-white dark:bg-slate-800 shadow-lg rounded-lg border border-slate-200 dark:border-slate-700">
          <h2 className="text-2xl font-semibold mb-4 text-slate-800 dark:text-slate-100 border-b dark:border-slate-600 pb-2">Overall Summary</h2>
          <div className="space-y-4">
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Key Developments</h3>
            <div className="space-y-3">
              {latestRun.overallSummary.keyDevelopments.map((keyDev, index) => (
                <div key={index} className="p-4 bg-slate-50 dark:bg-slate-700 rounded-md border-l-4 border-blue-500 dark:border-blue-400">
                  <p className="text-slate-700 dark:text-slate-300 mb-2">{keyDev.text}</p>
                  {keyDev.itemID && (
                    <div className="text-sm text-slate-500 dark:text-slate-400">
                      <span className="font-medium">Referenced Item ID:</span> {keyDev.itemID}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      <div className="mb-6">
        <h2 className="text-2xl font-semibold mb-4 text-slate-800 dark:text-slate-100">Entry Summaries</h2>
        {latestRun.entrySummaries && latestRun.entrySummaries.length > 0 ? (
          latestRun.entrySummaries.map((entry, index) => {
            return (
              <AccordionItem key={entry.results.id ? `${entry.results.id}-${index}`: index} title={entry.results.title || "Untitled Entry"}>
                <div className="space-y-3 p-2">
                  <div className="p-3 border border-slate-200 dark:border-slate-700 rounded-md bg-white dark:bg-slate-800 mb-3">
                    <h4 className="text-md font-semibold text-slate-700 dark:text-slate-300 mb-2">Details</h4>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-2 text-sm">
                      <div className="bg-slate-50 dark:bg-slate-700 p-2 rounded text-slate-800 dark:text-slate-200"><span className="font-medium text-slate-600 dark:text-slate-400">ID:</span> {entry.results.id || "N/A"}</div>
                      <div className="bg-slate-50 dark:bg-slate-700 p-2 rounded text-slate-800 dark:text-slate-200"><span className="font-medium text-slate-600 dark:text-slate-400">Relevant:</span> {entry.results.isRelevant ? 'Yes' : 'No'}</div>
                      <div className="bg-slate-50 dark:bg-slate-700 p-2 rounded text-slate-800 dark:text-slate-200"><span className="font-medium text-slate-600 dark:text-slate-400">Processing Time:</span> {typeof entry.processingTimeMs === 'number' ? entry.processingTimeMs : "N/A"} ms</div>
                      <div className="bg-slate-50 dark:bg-slate-700 p-2 rounded text-slate-800 dark:text-slate-200 col-span-full"><span className="font-medium text-slate-600 dark:text-slate-400">Link:</span> {entry.results.link || "N/A"}</div>
                    </div>
                  </div>
                  <TruncatableContent title="Raw Input" content={entry.rawInput || ""} />
                  <TruncatableContent title="Summary" content={entry.results.summary || ""} initialLineLimit={5} />
                  {entry.results.commentSummary && (
                      <TruncatableContent title="Comment Summary" content={entry.results.commentSummary || ""} initialLineLimit={5} />
                  )}
                  {entry.results.imageDescription && (  
                    <TruncatableContent title="Image Description" content={entry.results.imageDescription || ""} initialLineLimit={5} />
                  )}
                  {entry.results.webContentSummary && (  
                    <TruncatableContent title="Web Content Summary" content={entry.results.webContentSummary || ""} initialLineLimit={5} />
                  )}
                </div>
              </AccordionItem>
            );
          })
        ) : (
          <p className="text-slate-600 dark:text-slate-400">No entry summaries available for this run.</p>
        )}
      </div>
    </Layout>
  );
}
