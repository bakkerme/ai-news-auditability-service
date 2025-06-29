import Layout from '../../components/Layout';
import MetadataCard from '../../components/MetadataCard';
import AccordionItem from '../../components/AccordionItem';
import TruncatableContent from '../../components/TruncatableContent';
import BenchmarkControls from '../../components/BenchmarkControls';
import { RunData } from '../../types';

async function getRunData(runId: string): Promise<RunData> {
  try {
    const res = await fetch(`http://localhost:8080/v1/runs/${runId}`, {
      cache: 'no-store',
    });

    if (!res.ok) {
      const errorBody = await res.text();
      console.error(`Error fetching run data: ${res.status} ${res.statusText} - ${errorBody}`);
      throw new Error(`Failed to fetch run data: ${res.status} ${res.statusText}. Body: ${errorBody}`);
    }
    const data: RunData = await res.json();
    console.log('Run data:', data);
    return data;
  } catch (error) {
    console.error('Exception during fetch run data:', error);
    if (error instanceof Error && error.message.startsWith('Failed to fetch')) {
      throw error;
    }
    throw new Error(`An unexpected error occurred while fetching run data: ${error instanceof Error ? error.message : String(error)}`);
  }
}

interface RunPageProps {
  params: Promise<{
    runId: string;
  }>;
}

export default async function RunPage({ params }: RunPageProps) {
  const { runId } = await params;
  let runData: RunData | null = null;
  let fetchError: string | null = null;

  try {
    runData = await getRunData(runId);
  } catch (error) {
    console.error("Error in RunPage while fetching data:", error);
    fetchError = error instanceof Error ? error.message : "An unknown error occurred.";
  }

  if (fetchError || !runData) {
    return (
      <Layout>
        <div className="container mx-auto p-4">
          <h1 className="text-2xl font-bold text-red-600 dark:text-red-400 mb-4">Error Loading Run Data</h1>
          <p className="text-red-500 dark:text-red-300 text-center text-xl">
            Failed to load run data for ID: {runId}. Please check if the run exists or try again later.
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

  return (
    <Layout>
      <MetadataCard
        runId={runId}
        timestamp={runData.runDate}
        overallModel={runData.overallModelUsed}
        imageModel={runData.imageModelUsed}
        webContentModel={runData.webContentModelUsed}
        totalProcessingTime={runData.totalProcessingTime}
      />

      {/* Benchmark Section */}
      <div className="mb-6 p-6 bg-white dark:bg-slate-800 shadow-lg rounded-lg border border-slate-200 dark:border-slate-700">
        <h2 className="text-2xl font-semibold mb-4 text-slate-800 dark:text-slate-100 border-b dark:border-slate-600 pb-2">Benchmark</h2>
        <div className="flex flex-col sm:flex-row gap-4 items-center">
          <div className="flex-1">
            <p className="text-slate-600 dark:text-slate-400 mb-2">
              Evaluate the quality and relevance of AI-generated summaries for this run.
            </p>
            <p className="text-sm text-slate-500 dark:text-slate-500">
              Run ID: {runId}
            </p>
          </div>
          <BenchmarkControls runId={runId} />
        </div>
      </div>

      <div className="mb-6 p-6 bg-white dark:bg-slate-800 shadow-lg rounded-lg border border-slate-200 dark:border-slate-700">
        <h2 className="text-2xl font-semibold mb-4 text-slate-800 dark:text-slate-100 border-b dark:border-slate-600 pb-2">Persona Details</h2>
        <div className="space-y-4">
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Name</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">{runData.persona.name}</p>
          </div>
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Identity / System Prompt</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md whitespace-pre-wrap">{runData.persona.basePromptTask || "Not available"}</p>
          </div>
          <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Focus Areas</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">
              {Array.isArray(runData.persona.focusAreas) ? runData.persona.focusAreas.join(', ') : "Not available"}
            </p>
          </div>
           <div>
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Exclusion Criteria</h3>
            <p className="text-slate-600 dark:text-slate-400 p-2 bg-slate-50 dark:bg-slate-700 rounded-md">
              {Array.isArray(runData.persona.exclusionCriteria) ? runData.persona.exclusionCriteria.join(', ') : "Not available"}
            </p>
          </div>
        </div>
      </div>

      {/* Overall Summary Section */}
      {runData.overallSummary && runData.overallSummary.keyDevelopments && runData.overallSummary.keyDevelopments.length > 0 && (
        <div className="mb-6 p-6 bg-white dark:bg-slate-800 shadow-lg rounded-lg border border-slate-200 dark:border-slate-700">
          <h2 className="text-2xl font-semibold mb-4 text-slate-800 dark:text-slate-100 border-b dark:border-slate-600 pb-2">Overall Summary</h2>
          <div className="space-y-4">
            <h3 className="text-lg font-medium text-slate-700 dark:text-slate-300">Key Developments</h3>
            <div className="space-y-3">
              {runData.overallSummary.keyDevelopments.map((keyDev, index) => (
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
        {runData.entrySummaries && runData.entrySummaries.length > 0 ? (
          runData.entrySummaries.map((entry, index) => {
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
                  {entry.results.relevanceToCriteria && (
                      <TruncatableContent title="Relevance to Criteria" content={entry.results.relevanceToCriteria || ""} initialLineLimit={5} />
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