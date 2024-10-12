module key_storage_with_tamper (
    input wire clk,
    input wire reset,
    input wire tamper_detected,     // Tamper detection signal
    // Full key inputs
    input wire [255:0] rsa_full_key,
    input wire [255:0] ecdsa_full_key,
    input wire [255:0] eddsa_full_key,
    // Shard inputs
    input wire [255:0] rsa_shard,
    input wire [255:0] ecdsa_shard,
    input wire [255:0] eddsa_shard,
    // AES key for encryption and decryption
    input wire [255:0] aes_key,
    // Outputs
    output reg [255:0] rsa_full_key_out,
    output reg [255:0] rsa_shard_out,
    output reg [255:0] ecdsa_full_key_out,
    output reg [255:0] ecdsa_shard_out,
    output reg [255:0] eddsa_full_key_out,
    output reg [255:0] eddsa_shard_out,
    output reg [255:0] aes_key_out
);

    always @(posedge clk or posedge reset or posedge tamper_detected) begin
        if (reset || tamper_detected) begin
            // Zero out keys in case of tampering or reset
            rsa_full_key_out <= 256'b0;
            rsa_shard_out <= 256'b0;
            ecdsa_full_key_out <= 256'b0;
            ecdsa_shard_out <= 256'b0;
            eddsa_full_key_out <= 256'b0;
            eddsa_shard_out <= 256'b0;
            aes_key_out <= 256'b0;
        end else begin
            // Securely use keys
            rsa_full_key_out <= rsa_full_key;
            rsa_shard_out <= rsa_shard;
            ecdsa_full_key_out <= ecdsa_full_key;
            ecdsa_shard_out <= ecdsa_shard;
            eddsa_full_key_out <= eddsa_full_key;
            eddsa_shard_out <= eddsa_shard;
            aes_key_out <= aes_key;
        end
    end
endmodule