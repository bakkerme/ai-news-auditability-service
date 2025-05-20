# MVP UI Planning: Goal 1, Task 4 - Display basic UI for data

**Technology:**
- **Frontend Framework**: Next.js

**Detailed Subtask Breakdown:**

1.  **Design simple dashboard for viewing most recent run**
    *   **1.1. Define Data Requirements:**
        *   Identify all data points from the "most recent run" to be displayed. These will be sourced primarily from the `RunData` schema (see `#/components/schemas/RunData` in `api-doc.yaml`).
            *   **Run Identification:**
                *   Run ID: The identifier of the run being displayed (e.g., the `runId` from the API path `GET /runs/{runId}`).
            *   **Metadata:**
                *   Timestamp: `RunData.runDate`.
                *   Models Used:
                    *   Overall Model: `RunData.overallModelUsed`.
                    *   Image Model: `RunData.imageModelUsed` (if applicable and available).
                    *   Web Content Model: `RunData.webContentModelUsed` (if applicable and available).
                *   Execution Time: `RunData.totalProcessingTime` (overall) and/or `RunData.entrySummaries[0].processingTime` (for the first/primary entry).
                *   Token Counts: (e.g., prompt, completion, total) - Desirable for display, but not currently part of the `RunData` schema. This may require an API update or be omitted for the MVP.
            *   **Persona Details** (from `RunData.persona`, which refers to `#/components/schemas/Persona`):
                *   Name: `Persona.name`.
                *   Identity/System Prompt: `Persona.identity`.
                *   Focus Areas: `Persona.focusAreas`.
                *   (Note: `version` is not currently a field in the `Persona` schema).
            *   **Core Content** (assuming display of the first entry summary for MVP: `RunData.entrySummaries[0]`):
                *   Raw Input/Query: `EntrySummary.rawInput`.
                *   LLM Response (Summary): `EntrySummary.results.summary` (where `results` refers to `#/components/schemas/Item`).
    *   **1.2. Sketch UI Layout (Wireframe):**
        *   Header: Service Title (e.g., "AI News Auditability - Latest Run").
        *   Main Content Area:
            *   Section for "Run Metadata" (Timestamp, Model, Run ID).
            *   Section for "Prompt Breakdown":
                *   Display each prompt component clearly, perhaps in separate cards or collapsible sections (e.g., "System Prompt", "Persona Details").
            *   Section for "LLM Output": Display the full response from the LLM.
        *   Footer: (Optional - e.g., version info, links).
    *   **1.3. Choose a Styling Approach:**
        *   Decide on a CSS framework or methodology (e.g., Tailwind CSS for utility-first, Material-UI for pre-built components, CSS Modules for scoped styles).
2.  **Develop frontend components to display prompt components (Next.js)**
    *   **2.1. Setup Next.js Project:**
        *   Initialize a new Next.js application if one doesn't exist for this UI.
        *   Configure the chosen styling approach.
    *   **2.2. Create Main Page/View:**
        *   Create a page (e.g., `pages/index.js` or `pages/latest-run.js`) for the dashboard.
        *   Implement data fetching (e.g., using `getServerSideProps` or client-side fetching with `useEffect` and `fetch/axios` to call the `/api/v1/runs/latest` endpoint).
    *   **2.3. Create Reusable UI Components:**
        *   `Layout.js`: Main application layout (header, footer, main content area).
        *   `MetadataCard.js`: Component to display run metadata.
        *   `PromptComponentDisplay.js`: A generic component to display a named prompt segment (e.g., takes `title` and `content` props). Could be a card with a title and pre-formatted text.
        *   `LLMResponseDisplay.js`: Component to display the LLM's response, potentially with formatting.
    *   **2.4. Assemble the Dashboard:**
        *   Use the created components on the main page to render the fetched data according to the design.
        *   Ensure responsive design for different screen sizes if necessary, though "basic UI" might imply desktop-first.

3.  **Create basic navigation structure**
    *   **3.1. Implement Header Navigation:**
        *   Add a simple header to the `Layout.js` component.
        *   Initially, this might just be the service title.
        *   If other pages are planned (e.g., for "Goal 3: Auto Benchmarking Support"), placeholders or links can be added.
    *   **3.2. (Optional) Sidebar Navigation:**
        *   If the UI is expected to grow with multiple sections, consider a basic sidebar for navigation. For a single view of the "most recent run," this is likely not needed for the MVP. 