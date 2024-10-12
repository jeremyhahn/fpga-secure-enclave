module ed25519_signing_core (
    input wire clk,
    input wire reset,
    input wire [255:0] private_key,    // Full private key for Ed25519
    input wire [127:0] message_hash,   // Message hash
    input wire partial,                // Use partial signing (threshold)
    input wire [255:0] partial_key,    // Shard of the private key for partial signature
    output reg [255:0] signature_out,  // Ed25519 signature
    output reg done
);

    ed25519_full_signer ed25519_full (
        .clk(clk),
        .reset(reset),
        .private_key(private_key),
        .message_hash(message_hash),
        .signature(signature_out),
        .done(done)
    );

    ed25519_partial_signer ed25519_partial (
        .clk(clk),
        .reset(reset),
        .partial_key(partial_key),
        .message_hash(message_hash),
        .signature(signature_out),
        .done(done)
    );

    always @(posedge clk or posedge reset) begin
        if (reset) begin
            signature_out <= 256'b0;
            done <= 0;
        end else if (partial) begin
            // Use partial key signing if indicated
            ed25519_partial.partial_key <= partial_key;
            ed25519_partial.message_hash <= message_hash;
        end else begin
            // Use full key signing
            ed25519_full.private_key <= private_key;
            ed25519_full.message_hash <= message_hash;
        end
    end
endmodule