import React, { useState, useRef, useEffect } from "react";
import { jsPDF } from "jspdf";

const App: React.FC = () => {
  const [text, setText] = useState<string>("");
  const [inflSignature, setInflSignature] = useState<string>("");
  const [font, setFont] = useState<string>("Arial");
  const [title, setTitle] = useState<string>("");
  const [brand, setBrand] = useState<string>("");
  const [managerFullName, setManagerFullName] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [inflFullName, setInfluFullName] = useState<string>("");
  const [inflEmail, setInflEmail] = useState<string>("");

  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const inflCanvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (canvas) {
      const ctx = canvas.getContext("2d");
      if (ctx) {
        // Clear the canvas
        ctx.clearRect(0, 0, canvas.width, canvas.height);

        // Set the background color
        ctx.fillStyle = "white";
        ctx.fillRect(0, 0, canvas.width, canvas.height);

        // Set the font and text properties
        ctx.font = `bold 48px ${font}`;
        ctx.fillStyle = "black";
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";

        // Draw the text on the canvas
        ctx.fillText(text, canvas.width / 2, canvas.height / 2);
      }
    }
  }, [text, font]);

  useEffect(() => {
    const canvas = inflCanvasRef.current;
    if (canvas) {
      const ctx = canvas.getContext("2d");
      if (ctx) {
        // Clear the canvas
        ctx.clearRect(0, 0, canvas.width, canvas.height);

        // Set the background color
        ctx.fillStyle = "white";
        ctx.fillRect(0, 0, canvas.width, canvas.height);

        // Set the font and text properties
        ctx.font = `bold 48px ${font}`;
        ctx.fillStyle = "black";
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";

        // Draw the text on the canvas
        ctx.fillText(inflSignature, canvas.width / 2, canvas.height / 2);
      }
    }
  }, [inflSignature, font]);

  // Function to call the API and get the S3 upload URL
  const getS3UploadUrl = async (s3Key: string, userId: number) => {
    const response = await fetch("http://localhost:4000/contract", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        s3Key,
        userId,
      }),
    });

    if (!response.ok) {
      throw new Error("Failed to get S3 upload URL");
    }

    const data = await response.json();
    return data.url;
  };

  // Function to upload the PDF to S3 using the provided URL
  const uploadPdfToS3 = async (pdfData: Blob, uploadUrl: string) => {
    const response = await fetch(uploadUrl, {
      method: "PUT",
      body: pdfData,
    });

    if (!response.ok) {
      throw new Error("Failed to upload PDF to S3");
    }
  };

  // Function to upload the image to S3 using the provided URL
  const uploadImageToS3 = async (imgData: Blob, uploadUrl: string) => {
    const response = await fetch(uploadUrl, {
      method: "PUT",
      body: imgData,
    });

    if (!response.ok) {
      throw new Error("Failed to upload image to S3");
    }
  };

  const generateInflSignature = async () => {
    const canvas = inflCanvasRef.current;
    if (canvas) {
      const imgBlob = await new Promise<Blob>((resolve) => {
        canvas.toBlob((blob) => {
          resolve(blob as Blob);
        }, "image/png");
      });

      try {
        // Get the S3 upload URL from the API
        const s3Key = `influencer-signature-${inflSignature}.png`; // Replace with the desired S3 key
        const userId = 1; // Replace with the actual user ID
        const uploadUrl = await getS3UploadUrl(s3Key, userId);

        // Upload the image to S3
        await uploadImageToS3(imgBlob, uploadUrl);

        console.log("Influencer signature uploaded to S3 successfully");
      } catch (error) {
        console.error("Error uploading influencer signature to S3:", error);
      }
    }
  };

  const generatePDF = async () => {
    const canvas = canvasRef.current;
    if (canvas) {
      const imgData = canvas.toDataURL("image/png");

      const pdf = new jsPDF();

      // Set the initial y-coordinate for the text
      let y = 20;
      const lineHeight = 10;

      // Add contract text to the PDF
      pdf.text(`Title - ${title}`, 20, y);
      y += lineHeight;
      pdf.text(`Brand - ${brand}`, 20, y);
      y += lineHeight;
      pdf.text(`Manager Full Name - ${managerFullName}`, 20, y);
      y += lineHeight;
      pdf.text(`Email - ${email}`, 20, y);
      y += lineHeight;
      pdf.text(`Influencer Full Name - ${inflFullName}`, 20, y);
      y += lineHeight;
      pdf.text(`Influencer Email - ${inflEmail}`, 20, y);
      y += lineHeight;

      // Add the generated image to the left signature box
      pdf.addImage(imgData, "PNG", 20, y + 20, 80, 40);

      // Add the right signature box
      pdf.rect(120, y + 20, 80, 40);
      pdf.text("Signature", 140, y + 40);
      const pdfData = pdf.output("blob");

      try {
        // Get the S3 upload URL from the API
        const s3Key = `contract-${text}.pdf`; // Replace with the desired S3 key
        const userId = 1; // Replace with the actual user ID
        const uploadUrl = await getS3UploadUrl(s3Key, userId);

        // Upload the PDF to S3
        await uploadPdfToS3(pdfData, uploadUrl);

        console.log("PDF uploaded to S3 successfully");
      } catch (error) {
        console.error("Error uploading PDF to S3:", error);
      }
      // Save the PDF
    }
  };

  const handleTextChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setText(event.target.value);
  };

  const handleFontChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setFont(event.target.value);
  };

  // const handleDownload = () => {
  //   const canvas = canvasRef.current;
  //   if (canvas) {
  //     const dataURL = canvas.toDataURL("image/png");
  //     const link = document.createElement("a");
  //     link.href = dataURL;
  //     link.download = "generated-image.png";
  //     link.click();
  //   }
  // };

  return (
    <div>
      <div style={{ display: "flex" }}>
        <input
          type="text"
          value={text}
          onChange={handleTextChange}
          placeholder="Enter signature"
        />
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Enter title"
        />
        <input
          type="text"
          value={brand}
          onChange={(e) => setBrand(e.target.value)}
          placeholder="Enter brand"
        />
        <input
          type="text"
          value={managerFullName}
          onChange={(e) => setManagerFullName(e.target.value)}
          placeholder="Enter manager full name"
        />
        <input
          type="text"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Enter manager email"
        />
        <input
          type="text"
          value={inflFullName}
          onChange={(e) => setInfluFullName(e.target.value)}
          placeholder="Enter influencer full name"
        />
        <input
          type="text"
          value={inflEmail}
          onChange={(e) => setInflEmail(e.target.value)}
          placeholder="Enter influencer email"
        />
      </div>
      <select value={font} onChange={handleFontChange}>
        <option value="Arial">Arial</option>
        <option value="Verdana">Verdana</option>
        <option value="Times New Roman">Times New Roman</option>
      </select>
      {/* <button onClick={handleDownload}>Download Image</button>*/}
      <canvas ref={canvasRef} width={400} height={200} />
      <button onClick={generatePDF}>Generate PDF</button>
      <div style={{ marginTop: 200 }}>
        <canvas ref={inflCanvasRef} width={400} height={200} />
        <button onClick={generateInflSignature}>Generate Infl signature</button>
        <input
          type="text"
          value={inflSignature}
          onChange={(e) => setInflSignature(e.target.value)}
          placeholder="Enter signature"
        />
      </div>
    </div>
  );
};

export default App;
