import React, { useState } from 'react';
import { Flag, Youtube, Headphones, FileText, Database, BookOpen, ExternalLink } from 'lucide-react';

function App() {
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [backendUrl, setBackendUrl] = useState('');

  const startBackendServer = () => {
    setIsLoading(true);
    setMessage('Starting the backend application...');
    
    // Since we can't directly start the backend from the browser in this environment,
    // we'll provide instructions to the user
    setMessage('Please run "npm run start-backend" in a terminal window to start the backend server. This will install the required Python dependencies using UV and start the Gradio application.');
    
    // Set a default URL for development purposes
    setBackendUrl('http://localhost:7860');
    setIsLoading(false);
  };

  const openBackendApp = () => {
    if (backendUrl) {
      window.open(backendUrl, '_blank');
    } else {
      startBackendServer();
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header with Italian flag colors */}
      <header className="bg-gradient-to-r from-green-600 via-white to-red-600 p-6 shadow-md">
        <div className="max-w-7xl mx-auto">
          <h1 className="text-3xl font-bold text-center text-gray-800">
            🇮🇹 Italian Language Learning Platform 🇮🇹
          </h1>
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
        <div className="bg-white rounded-lg shadow-xl overflow-hidden">
          {/* Hero section */}
          <div className="relative">
            <img 
              src="https://images.unsplash.com/photo-1516483638261-f4dbaf036963?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1200&q=80" 
              alt="Italian landscape" 
              className="w-full h-96 object-cover"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent flex items-end">
              <div className="p-8 text-white">
                <h2 className="text-4xl font-bold mb-4">Improve Your Italian Listening Skills</h2>
                <p className="text-xl mb-6">
                  Practice with authentic content and AI-powered exercises
                </p>
                <button 
                  onClick={openBackendApp}
                  disabled={isLoading}
                  className="bg-green-600 hover:bg-green-700 text-white font-bold py-3 px-6 rounded-lg flex items-center transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isLoading ? (
                    <span className="flex items-center">
                      <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Starting...
                    </span>
                  ) : backendUrl ? (
                    <span className="flex items-center">
                      Open Application <ExternalLink className="ml-2 h-5 w-5" />
                    </span>
                  ) : (
                    <span className="flex items-center">
                      Launch Application <ExternalLink className="ml-2 h-5 w-5" />
                    </span>
                  )}
                </button>
                {message && (
                  <div className="mt-4 text-sm bg-white/10 p-3 rounded">
                    <p>{message}</p>
                    {backendUrl && (
                      <p className="mt-2">
                        Backend URL: <code className="bg-black/20 px-2 py-1 rounded">{backendUrl}</code>
                      </p>
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Features section */}
          <div className="p-8">
            <h3 className="text-2xl font-bold text-gray-800 mb-6">Key Features</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <FeatureCard 
                icon={<Youtube className="h-8 w-8 text-red-500" />}
                title="YouTube Transcripts"
                description="Extract transcripts from Italian YouTube videos for learning materials."
              />
              
              <FeatureCard 
                icon={<Headphones className="h-8 w-8 text-purple-500" />}
                title="Whisper Transcription"
                description="Generate accurate transcripts from audio using OpenAI's Whisper model."
              />
              
              <FeatureCard 
                icon={<FileText className="h-8 w-8 text-blue-500" />}
                title="OCR Extraction"
                description="Extract text from video frames using specialized Italian OCR."
              />
              
              <FeatureCard 
                icon={<Database className="h-8 w-8 text-green-500" />}
                title="RAG Implementation"
                description="Use retrieval-augmented generation for contextual learning."
              />
              
              <FeatureCard 
                icon={<BookOpen className="h-8 w-8 text-amber-500" />}
                title="Interactive Learning"
                description="Practice with AI-generated exercises tailored to your level."
              />
              
              <FeatureCard 
                icon={<Flag className="h-8 w-8 text-green-600" />}
                title="Italian-Focused"
                description="Specialized tools and models optimized for the Italian language."
              />
            </div>
          </div>

          {/* How it works section */}
          <div className="bg-gray-50 p-8">
            <h3 className="text-2xl font-bold text-gray-800 mb-6">How It Works</h3>
            
            <div className="space-y-4">
              <Step 
                number={1} 
                title="Extract Content" 
                description="Download transcripts from Italian YouTube videos or generate them using Whisper."
              />
              
              <Step 
                number={2} 
                title="Process and Store" 
                description="Extract text using OCR and store in a vector database for efficient retrieval."
              />
              
              <Step 
                number={3} 
                title="Generate Exercises" 
                description="Create personalized listening comprehension exercises based on the content."
              />
              
              <Step 
                number={4} 
                title="Practice and Learn" 
                description="Interact with the exercises and get feedback to improve your Italian skills."
              />
            </div>
          </div>

          {/* Setup Instructions */}
          <div className="p-8 border-t border-gray-200">
            <h3 className="text-2xl font-bold text-gray-800 mb-6">Setup Instructions</h3>
            
            <div className="bg-gray-50 p-6 rounded-lg border border-gray-200">
              <h4 className="text-lg font-semibold mb-4">Before You Begin</h4>
              
              <ol className="list-decimal pl-5 space-y-3">
                <li>
                  <p className="font-medium">Configure your environment variables</p>
                  <p className="text-gray-600 mt-1">Copy <code className="bg-gray-200 px-2 py-1 rounded">.env.example</code> to <code className="bg-gray-200 px-2 py-1 rounded">.env</code> and fill in your API keys.</p>
                </li>
                <li>
                  <p className="font-medium">Install dependencies</p>
                  <p className="text-gray-600 mt-1">Make sure all required packages are installed:</p>
                  <div className="bg-gray-800 text-white p-3 rounded mt-2 overflow-x-auto">
                    <code>npm install</code>
                  </div>
                </li>
                <li>
                  <p className="font-medium">Start the backend server</p>
                  <p className="text-gray-600 mt-1">Run the following command in your terminal:</p>
                  <div className="bg-gray-800 text-white p-3 rounded mt-2 overflow-x-auto">
                    <code>npm run start-backend</code>
                  </div>
                  <p className="text-gray-600 mt-1">This will automatically install the required Python dependencies using UV and start the Gradio application.</p>
                </li>
                <li>
                  <p className="font-medium">Access the application</p>
                  <p className="text-gray-600 mt-1">Once the backend is running, click the "Launch Application" button above or navigate to <a href="http://localhost:7860" target="_blank" className="text-green-600 hover:underline">http://localhost:7860</a> in your browser.</p>
                </li>
              </ol>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-gray-800 text-white py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="mb-4 md:mb-0">
              <h4 className="text-xl font-bold">Italian Language Learning Platform</h4>
              <p className="text-gray-400">Powered by AI and authentic content</p>
            </div>
            <div className="flex space-x-4">
              <FooterLink href="#" text="About" />
              <FooterLink href="#" text="Documentation" />
              <FooterLink href="#" text="GitHub" />
              <FooterLink href="#" text="Contact" />
            </div>
          </div>
          <div className="mt-8 pt-8 border-t border-gray-700 text-center text-gray-400">
            <p>© 2025 Italian Language Learning Platform. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({ icon, title, description }) {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md border border-gray-100 hover:shadow-lg transition-shadow">
      <div className="mb-4">{icon}</div>
      <h4 className="text-xl font-semibold mb-2 text-gray-800">{title}</h4>
      <p className="text-gray-600">{description}</p>
    </div>
  );
}

function Step({ number, title, description }) {
  return (
    <div className="flex items-start">
      <div className="flex-shrink-0 bg-green-600 text-white rounded-full w-10 h-10 flex items-center justify-center font-bold text-lg">
        {number}
      </div>
      <div className="ml-4">
        <h4 className="text-lg font-semibold text-gray-800">{title}</h4>
        <p className="text-gray-600">{description}</p>
      </div>
    </div>
  );
}

function FooterLink({ href, text }) {
  return (
    <a 
      href={href} 
      className="text-gray-300 hover:text-white transition-colors"
    >
      {text}
    </a>
  );
}

export default App;