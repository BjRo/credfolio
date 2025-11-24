import { DefaultService } from "./generated";

export const uploadReferenceLetter = async (file: File) => {
	// The generated client expects a Blob, File inherits from Blob so it's fine.
	// The type definition says { file?: Blob }, so we pass it like that.
	return DefaultService.uploadReferenceLetter({ file });
};
