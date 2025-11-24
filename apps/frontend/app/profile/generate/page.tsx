"use client";

import { useState } from "react";
import ReferenceLetterUpload from "../../../components/profile/ReferenceLetterUpload";
import GenerateProfileButton from "../../../components/profile/GenerateProfileButton";
import type { Profile } from "../../../lib/api/generated/models/Profile";

export default function GeneratePage() {
  const [profile, setProfile] = useState<Profile | null>(null);

  const handleGenerateComplete = (profileData: Profile) => {
    setProfile(profileData);
    // Redirect to profile view or show success
    console.log("Generated Profile:", profileData);
  };

  return (
    <div className="container mx-auto p-8 max-w-2xl">
      <h1 className="text-3xl font-bold mb-6">Generate Smart Profile</h1>

      <div className="mb-8">
        <h2 className="text-xl font-semibold mb-4">
          1. Upload Reference Letters
        </h2>
        <div className="space-y-4">
          <ReferenceLetterUpload />
        </div>
        <p className="text-sm text-gray-500 mt-2">
          You can upload multiple letters one by one.
        </p>
      </div>

      <div className="mb-8">
        <h2 className="text-xl font-semibold mb-4">2. Generate Profile</h2>
        <p className="mb-4 text-gray-600">
          Our AI will analyze your reference letters to extract your work
          experience, skills, and credibility highlights.
        </p>
        <GenerateProfileButton onGenerateComplete={handleGenerateComplete} />
      </div>

      {profile && (
        <div className="bg-green-50 border border-green-200 rounded p-4">
          <h3 className="text-green-800 font-medium">Success!</h3>
          <p className="text-green-700">
            Profile generated with {profile.workExperiences?.length} work
            experiences.
          </p>
          <pre className="bg-white p-2 mt-2 rounded text-xs overflow-auto max-h-60">
            {JSON.stringify(profile, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}
